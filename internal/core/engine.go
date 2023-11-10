package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anacrolix/chansync"

	"github.com/varbhat/exatorrent/internal/db"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type Eng struct {
	Tconfig TorConfig
	Econfig EngConfig
	PsqlCon string

	Torc *torrent.Client

	onCloseMap map[metainfo.Hash]*chansync.Flag

	PcDb      db.PcDb
	TorDb     db.TorrentDb
	TrackerDB db.TrackerDb
	FsDb      db.FileStateDb
	LsDb      db.LockStateDb
	UDb       db.UserDb
	TUDb      db.TorrentUserDb
}

// AddfromSpec Adds Torrent by Torrent Spec
func AddfromSpec(User string, spec *torrent.TorrentSpec, dontstart bool, nofsdb bool) {
	if spec == nil {
		return
	}

	if Engine.Econfig.GetDTU() {
		RemTrackersSpec(spec)
	}

	Info.Println("Adding Torrent")
	var trnt *torrent.Torrent
	var new bool
	var err error

	// Add Torrent to Bittorrent Client
	if Engine.Torc != nil {
		trnt, new, err = Engine.Torc.AddTorrentSpec(spec)
		if err != nil {
			Warn.Println("Error adding Torrent Spec", err)
			MainHub.SendMsgU(User, "nfn", trnt.InfoHash().HexString(), "error", "Error adding Torrent Spec")
			return
		}

		if !new {
			Info.Printf("Torrent %s is not new", spec.InfoHash)
			MainHub.SendMsgU(User, "nfn", trnt.InfoHash().HexString(), "warning", "Torrent is not new")
			if User != "" {
				_ = Engine.TUDb.Add(User, trnt.InfoHash())
			}
			return
		}
		if Engine.Econfig.GetLBD() {
			_ = Engine.LsDb.Lock(spec.InfoHash)
			Info.Println("lock of torrent ", spec.InfoHash, " is set to true (due to config option)")
		}
	} else {
		Err.Println("Torrent Client is nil")
		return
	}

	if User != "" {
		_ = Engine.TUDb.Add(User, trnt.InfoHash())
		Info.Println("Torrent: ", spec.InfoHash, " User: ", User)
	}

	// Add Trackers to Torrent
	count := Engine.TrackerDB.Count()
	if count > 0 && (len(Engine.Econfig.TrackerListURLs) != 0 || len(trnt.Metainfo().AnnounceList) == 0) {
		trnt.AddTrackers([][]string{Engine.TrackerDB.Get()})
		Info.Println("Added ", count, " Trackers for Added Torrent")
	}

	ih := trnt.InfoHash()
	hasadded := Engine.TorDb.Exists(ih)
	if !hasadded {
		_ = Engine.TorDb.Add(ih)
		Info.Println("Added Torrent to Torrent Database")
	}

	var cancelled Mutbool

	go func() {
		defer func() {
			if err := recover(); err != nil {
				Warn.Println(err)
			}
		}()

		var notmerged bool = true
		if !Engine.Econfig.GetDLC() {
			Info.Println("Getting Metainfo from Local Cache")
			sp, err := GetMetaCache(ih)
			if err != nil {
				Info.Println("Torrent Metainfo is not present in Cache: ", err)
				notmerged = true
			} else {

				if !(cancelled.Get()) {
					if Engine.Econfig.GetDTC() {
						RemTrackersSpec(sp)
					}
					if sp != nil {
						err := trnt.MergeSpec(sp)
						if err != nil {
							Warn.Println(err)
							notmerged = true
						} else {
							Info.Println("Torrent Metainfo merged from Local Cache")
							MainHub.SendMsgU(User, "nfn", ih.HexString(), "info", "Metainfo of Torrent "+ih.HexString()+" merged from Local Cache")
							notmerged = false
						}
					}
				}
			}
		}

		if len(Engine.Econfig.GetOCU()) != 0 && notmerged {
			Info.Println("Torrent Metainfo requested from Online Cache")
			uspec, err := SpecfromURL(fmt.Sprintf(Engine.Econfig.GetOCU(), strings.TrimSpace(trnt.InfoHash().HexString())))
			if err != nil {
				Warn.Println(err)
			} else {
				if !(cancelled.Get()) {
					if Engine.Econfig.GetDTC() {
						RemTrackersSpec(uspec)
					}
					if uspec != nil {
						trmerr := trnt.MergeSpec(uspec)
						if trmerr != nil {
							Warn.Println(err)
							return
						}
						Info.Println("Torrent Metainfo merged form Online Cache")
						MainHub.SendMsgU(User, "nfn", ih.HexString(), "info", "Metainfo of Torrent merged from Online Cache")
					}
				}
			}
		}

	}()

	select {
	case <-Engine.Torc.Closed():
		cancelled.Set(true)
		Warn.Println("Torrent Client Closed! So,Torrent Not Started")
		MainHub.SendMsgU(User, "nfn", ih.HexString(), "error", "Torrent Client Closed! So,Torrent Not Added")
	case <-trnt.Closed():
		cancelled.Set(true)
		Warn.Println("Torrent Task that was being added was Deleted !")
		MainHub.SendMsgU(User, "nfn", ih.HexString(), "warning", "Torrent Task		cancelled = true that was being added was Deleted !")
	case <-trnt.GotInfo():
		Info.Println("Torrent ", ih, " is Loaded")
		MainHub.SendMsgU(User, "nfn", ih.HexString(), "info", "Torrent is Loaded")
		if !dontstart {
			go StartTorrent(User, ih, nofsdb)
		}
		if !Engine.Econfig.GetDLC() {
			AddMetaCache(ih, trnt.Metainfo())
		}
	}
}

// StartTorrent Starts Torrent given infohash
func StartTorrent(User string, infohash metainfo.Hash, nofsdb bool) {
	if User != "" {
		if !Engine.TUDb.HasUser(User, infohash.HexString()) {
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	}

	Info.Println("Starting Torrent for Infohash ", infohash)

	t, err := Engine.TorDb.GetTorrent(infohash)
	if err != nil {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
		Warn.Println(err)
		return
	}
	if t.Started {
		Info.Println("Already Started")
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "info", "Torrent is Already Started")
	}

	trnt, ok := Engine.Torc.Torrent(infohash)
	if !ok {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
		Warn.Println("Error fetching torrent of infohash ", infohash, " from the client")
		return
	}

	if trnt.Info() != nil {
		trnt.DownloadAll()
	} else {
		Warn.Println("Torrent couldn't be Started Because Metainfo is not yet loaded")
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent couldn't be Started Because Metainfo is not yet loaded")
		return
	}
	err = Engine.TorDb.Start(infohash)
	if err != nil {
		Warn.Println(err)
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent couldn't be Started")
		return
	}

	if !nofsdb {
		err = Engine.FsDb.Delete(infohash)
		if err != nil {
			Warn.Println(err)
		}
	}

	Info.Println("Torrent ", infohash, " Started by ", User)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "Torrent Started")

	if Engine.Econfig.GetListenC() {
		go func() {
			if _, ok := Engine.onCloseMap[infohash]; !ok {
				Engine.onCloseMap[infohash] = &trnt.Complete
				Info.Println("Listening for Completion of Torrent ", infohash)
				<-trnt.Complete.On()
				delete(Engine.onCloseMap, infohash)

				_, err := Engine.TorDb.GetTorrent(infohash)
				if err != nil {
					Info.Println(infohash, " Removed")
				} else {
					Info.Println(infohash, " Completed")
					if Engine.Econfig.GetNOC() {
						MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "Torrent Completed")
					}
					hpu := Engine.Econfig.GetHPU()
					if hpu != "" {
						trntname := ""
						if trnt != nil {
							trntname = trnt.Name()
						}
						sendPostReq(infohash, hpu, trntname)
					}
				}
			}
		}()
	}
}

// StopTorrent Stops Torrent given infohash
func StopTorrent(User string, infohash metainfo.Hash) {
	if User != "" {
		if !Engine.TUDb.HasUser(User, infohash.HexString()) {
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	}

	Warn.Println("Stopping Torrent for Infohash ", infohash)

	t, err := Engine.TorDb.GetTorrent(infohash)
	if err != nil {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
		Warn.Println(err)
		return
	}
	if !t.Started {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "info", "Torrent Already Stopped !")
		Warn.Println("Already stopped")
	}

	trnt, ok := Engine.Torc.Torrent(infohash)
	if !ok {
		Warn.Println("Error fetching torrent of infohash ", infohash, " from the client")
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
		return
	}

	if trnt.Info() != nil {
		trnt.CancelPieces(0, trnt.NumPieces())
	} else {
		Warn.Println("Torrent can't be Stopped Because Metainfo is not yet received")
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent couldn't be Stopped Because Metainfo is not yet received")
		return
	}

	err = Engine.TorDb.SetStarted(infohash, false)
	if err != nil {
		Warn.Println(err)
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent couldn't be Stopped")
		return
	}

	err = Engine.FsDb.Delete(infohash)
	if err != nil {
		Warn.Println(err)
	}

	Info.Println("Torrent ", infohash, " Stopped by ", User)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "Torrent Stopped")
}

// RemoveTorrent Removes Torrent from Torrent Client
func RemoveTorrent(User string, infohash metainfo.Hash) {
	if User != "" {
		if !Engine.TUDb.HasUser(User, infohash.HexString()) {
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	}

	StopTorrent(User, infohash)
	Warn.Println("Removing Torrent of Infohash ", infohash)

	trnt, ok := Engine.Torc.Torrent(infohash)
	if !ok {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
		Warn.Println("Error fetching torrent of infohash ", infohash, " from the client")
		return
	}

	trnt.Drop()

	if Engine.Econfig.GetListenC() {
		trnt.Complete.Set()
		trnt.Complete.Clear()
	}

	err := Engine.TorDb.Delete(infohash)
	if err != nil {
		Warn.Println(err)
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent couldn't be Removed")
	}

	err = Engine.FsDb.Delete(infohash)
	if err != nil {
		Warn.Println(err)
	}
	Info.Println("Torrent ", infohash, " Removed by ", User)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "Torrent Removed")
}

// DeleteTorrent in addition to Removing Torrent from Client also Deletes it from Storage
func DeleteTorrent(User string, infohash metainfo.Hash) {
	if User != "" {
		if !Engine.TUDb.HasUser(User, infohash.HexString()) {
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	}

	RemoveTorrent(User, infohash)

	var err error
	if infohash.HexString() != "" {
		err = os.RemoveAll(filepath.Join(Dirconfig.TrntDir, infohash.HexString()))
		if err != nil {
			Warn.Println(err)
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	} else {
		Warn.Println("Unknown Error")
		return
	}

	Engine.PcDb.Delete(infohash)

	err = Engine.TUDb.RemoveAllMi(infohash)
	if err != nil {
		Warn.Println(err)
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent couldn't be Deleted")
	}

	_ = Engine.LsDb.Unlock(infohash)

	if !Engine.Econfig.DRCI() {
		RemMetaCache(infohash)
	}
	Info.Println("Torrent ", infohash, " Deleted by ", User)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "Torrent Deleted")
}

// StartFile Starts File in Torrent given infohash and Filepath
func StartFile(User string, infohash metainfo.Hash, fp string) {
	fp = filepath.ToSlash(fp)
	if fp == "" {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Filepath Invalid!")
		return
	}
	if User != "" {
		if !Engine.TUDb.HasUser(User, infohash.HexString()) {
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	}
	Info.Println("Starting File in Torrent for Infohash ", infohash, "and Filepath ", fp)

	t, ok := Engine.Torc.Torrent(infohash)
	if !ok {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
		Warn.Println("Error fetching torrent of infohash ", infohash, " from the client")
		return
	}

	// Get File from Given Torrent
	var f *torrent.File
	if t.Info() != nil {
		for _, file := range t.Files() {
			if file.Path() == fp {
				f = file
				break
			}
		}
	} else {
		Warn.Println("File ", fp, "can't be Started Because Metainfo is not yet received")
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "File "+fp+" couldn't be Started Because Metainfo is not yet received")
		return
	}

	// If File is missing , return the error
	if f == nil {
		Warn.Printf("%s not present", fp)
		MainHub.SendMsgU(User, "error", infohash.HexString(), "info", "File"+fp+" not present")
		return
	}

	// Set Priority to Normal / Start the Torrent
	f.SetPriority(torrent.PiecePriorityNormal)

	var err error
	err = Engine.FsDb.Deletefile(fp, infohash)
	if err != nil {
		Warn.Println(err)
	}

	if !Engine.TorDb.HasStarted(infohash.HexString()) {
		err = Engine.FsDb.Add(f.Path(), infohash)
		if err != nil {
			Warn.Println(err)
		}

	}

	Info.Println("File", fp, "of Torrent", infohash, " Started by ", User)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "File "+fp+" started")

}

// StopFile Stops File in Torrent given infohash and Filepath
func StopFile(User string, infohash metainfo.Hash, fp string) {
	fp = filepath.ToSlash(fp)
	if fp == "" {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Filepath Invalid!")
		return
	}
	if User != "" {
		if !Engine.TUDb.HasUser(User, infohash.HexString()) {
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	}

	Warn.Println("Stopping File in Torrent for Infohash ", infohash, "and Filepath ", fp)

	t, ok := Engine.Torc.Torrent(infohash)
	if !ok {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
		Warn.Println("Error fetching torrent of infohash ", infohash, " from the client")
		return
	}

	// Get File from Given Torrent
	var f *torrent.File
	if t.Info() != nil {
		for _, file := range t.Files() {
			if file.Path() == fp {
				f = file
				break
			}
		}
	} else {
		Warn.Println("File ", fp, "can't be Stopped Because Metainfo is not yet received")
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "File "+fp+" couldn't be Stopped Because Metainfo is not yet received")
		return
	}

	// If File is missing , return the error
	if f == nil {
		Warn.Printf("%s not present", fp)
		MainHub.SendMsgU(User, infohash.HexString(), "error", "info", "File"+fp+" not present")
		return
	}

	// Set Priority to None ( Stop the File)
	f.SetPriority(torrent.PiecePriorityNone)

	var err error
	err = Engine.FsDb.Deletefile(fp, infohash)
	if err != nil {
		Warn.Println(err)
	}

	if Engine.TorDb.HasStarted(infohash.HexString()) {
		err = Engine.FsDb.Add(f.Path(), infohash)
		if err != nil {
			Warn.Println(err)
		}
	}

	Info.Println("File", fp, "of Torrent", infohash, " Stopped by ", User)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "File "+fp+" stopped")
}

// DeleteFilePath deletes the file or Folder
func DeleteFilePath(User string, infohash metainfo.Hash, fp string) {
	if fp == "" {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Filepath Invalid!")
		return
	}
	if infohash.HexString() == "" {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Infohash Invalid!")
		return
	}

	fp = filepath.ToSlash(fp)
	if User != "" {
		if !Engine.TUDb.HasUser(User, infohash.HexString()) {
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	}

	Warn.Println("Deleting File in Torrent for Infohash ", infohash, "and Filepath ", fp)

	_, ok := Engine.Torc.Torrent(infohash)
	if ok {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Currently Present in Client. Please Remove it from Torrent Client First!")
		return
	}

	file := filepath.Join(Dirconfig.TrntDir, infohash.HexString(), filepath.FromSlash(fp))
	if !strings.HasPrefix(file, filepath.Join(Dirconfig.TrntDir, infohash.HexString())) {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "FilePath Invalid!")
		return
	}
	if file == filepath.Join(Dirconfig.TrntDir, infohash.HexString()) {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Cannot Delete whole Torrent Directory")
		return
	}

	err := os.RemoveAll(file)
	if err != nil {
		Warn.Println("DeleteFile error by Username: ", User, "Filepath: ", fp, " ", err)
		MainHub.SendMsgU(User, "error", infohash.HexString(), "info", "Error removing File "+fp)
		return
	}

	Info.Println("File", fp, "of Torrent", infohash, " Deleted by ", User)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "File"+fp+" Deleted")
}

// AbandonTorrent unlinks Torrent from User
func AbandonTorrent(User string, infohash metainfo.Hash) {

	err := Engine.TUDb.Remove(User, infohash)
	if err != nil {
		Warn.Println("AbandonTorrent error", err)
		MainHub.SendMsgU(User, "error", infohash.HexString(), "info", "AbandonTorrent error")
		return
	}

	Info.Println("User ", User, "  abandoned ", infohash)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "Torrent Abandoned")
}

// AddTrackerstoTorrent Adds Tracker to Torrent
func AddTrackerstoTorrent(User string, infohash metainfo.Hash, announcelist [][]string) {
	if User != "" {
		if !Engine.TUDb.HasUser(User, infohash.HexString()) {
			MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
			return
		}
	}

	trnt, ok := Engine.Torc.Torrent(infohash)
	if !ok {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Torrent Not Present!")
		Warn.Println("Error fetching torrent of infohash ", infohash, " from the client")
		return
	}

	if !(Engine.Econfig.GetDTU()) {
		trnt.AddTrackers(announcelist)
	} else {
		MainHub.SendMsgU(User, "nfn", infohash.HexString(), "error", "Trackers Couldn't be Added!")
		return
	}

	Info.Println("Trackers added to torrent ", infohash, "by ", User)
	MainHub.SendMsgU(User, "resp", infohash.HexString(), "success", "Trackers added")
}
