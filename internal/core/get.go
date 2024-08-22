package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type Torrent1 struct {
	Infohash         string   `json:"infohash"`
	Name             string   `json:"name,omitempty"`
	BytesCompleted   int64    `json:"bytescompleted,omitempty"`
	BytesMissing     int64    `json:"bytesmissing,omitempty"`
	BytesWritten     int64    `json:"byteswritten,omitempty"`
	Length           int64    `json:"length,omitempty"`
	State            string   `json:"state"`
	Seeding          bool     `json:"seeding,omitempty"`
	Private          bool     `json:"private,omitempty"`
	CreationDate     int64    `json:"creationdate,omitempty"`
	AddedDate        int64    `json:"addeddate,omitempty"`
	StartedDate      int64    `json:"starteddate,omitempty"`
	TotalPeers       int      `json:"totalpeers,omitempty"`
	ActivePeers      int      `json:"activepeers,omitempty"`
	ConnectedSeeders int      `json:"connectedseeders,omitempty"`
	AnnounceList     []string `json:"announcelist,omitempty"`
}

type Torrent2 struct {
	Torrent1
	StartedAt time.Time `json:"startedat"`
	AddedAt   time.Time `json:"addedat"`
}

type FileInfo struct {
	BytesCompleted int64  `json:"bytescompleted"`
	DisplayPath    string `json:"displaypath"`
	Length         int64  `json:"length"`
	Offset         int64  `json:"offset"`
	Path           string `json:"path"`
	Priority       byte   `json:"priority"`
}

type FsFileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isdir"`
}

type PeerInfo struct {
	RemoteAddr           string  `json:"remoteaddr"`
	PeerClientName       string  `json:"peerclientname"`
	DownloadRate         float64 `json:"downloadrate"`
	PeerPreferEncryption bool    `json:"peerpreferencryption"`
}

func GetTorrent(ih metainfo.Hash) (ret Torrent1) {
	ret.Infohash = ih.HexString()
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		ret.State = "removed"
		return
	}
	ret.Name = t.Name()
	if t == nil || t.Info() == nil {
		ret.State = "loading"
		return
	}
	if Engine.TorDb.HasStarted(ih.HexString()) {
		ret.State = "active"
	} else {
		ret.State = "inactive"
	}
	tdb, _ := Engine.TorDb.GetTorrent(t.InfoHash())
	if tdb != nil {
		ret.AddedDate = tdb.AddedAt.Unix()
		ret.StartedDate = tdb.StartedAt.Unix()
	}

	stats := t.Stats()
	info := t.Metainfo()

	if t.Info().Private != nil {
		ret.Private = *t.Info().Private
	}

	ret.Length = t.Length()
	ret.BytesCompleted = t.BytesCompleted()
	ret.BytesMissing = t.BytesMissing()
	ret.BytesWritten = stats.BytesWrittenData.Int64()
	ret.Seeding = t.Seeding()
	ret.CreationDate = info.CreationDate
	ret.TotalPeers = stats.TotalPeers
	ret.ActivePeers = stats.ActivePeers
	ret.ConnectedSeeders = stats.ConnectedSeeders
	ret.AnnounceList = info.AnnounceList.DistinctValues()

	return
}

func GetTorrents(lt []metainfo.Hash) (ret []byte) {
	var tinits []*Torrent1
	for _, ih := range lt {
		tinit := GetTorrent(ih)
		tinits = append(tinits, &tinit)
	}
	ret, _ = json.Marshal(DataMsg{Type: "torrentstream", Data: tinits})
	return
}

func GetTorrentInfo(ih metainfo.Hash) (ret []byte) {
	tinit := GetTorrent(ih)
	ret, _ = json.Marshal(DataMsg{Type: "torrentinfo", Data: tinit})
	return
}

func GetTorrentInfoStat(ih metainfo.Hash) (ret []byte) {
	trnt, err := Engine.TorDb.GetTorrent(ih)
	if err == nil {
		ret, _ = json.Marshal(DataMsg{Type: "torrentinfostat", Data: trnt})
	}
	return
}

func GetTorrentStats(ih metainfo.Hash) (ret []byte) {
	defer func() {
		if r := recover(); r != nil {
			Warn.Println("There was Panic in GetTorrentStats")
		}
	}()
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		ret, _ = json.Marshal(DataMsg{Type: "torrentstats", Data: nil})
		return
	}
	if t == nil || t.Info() == nil {
		ret, _ = json.Marshal(DataMsg{Type: "torrentstats", Data: nil})
		return
	}
	ts := t.Stats()
	ret, _ = json.Marshal(DataMsg{Data: &ts, Infohash: ih.HexString(), Type: "torrentstats"})
	return
}

func GetFsDirInfo(ih metainfo.Hash, fp string) (ret []byte) {
	fp = filepath.ToSlash(fp)
	ret, _ = json.Marshal(DataMsg{Type: "fsdirinfo", Data: nil})
	ihs := ih.HexString()
	if ihs == "" {
		Warn.Println("empty infohash")
		return
	}
	dirpath := filepath.Join(Dirconfig.TrntDir, ihs, fp)

	if !strings.HasPrefix(dirpath, filepath.Join(Dirconfig.TrntDir, ihs)) {
		return
	}
	rl, err := os.ReadDir(filepath.FromSlash(dirpath))
	if err != nil {
		Warn.Println(err.Error())
		return
	}

	var retfiles []FsFileInfo
	for _, eachdirentry := range rl {
		var retfile FsFileInfo
		inf, err := eachdirentry.Info()
		if err != nil {
			continue
		}
		retfile.Name = inf.Name()
		retfile.IsDir = inf.IsDir()
		retfile.Size = inf.Size()
		retfile.Path = filepath.ToSlash(filepath.Join(fp, retfile.Name))
		retfiles = append(retfiles, retfile)
	}
	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "fsdirinfo", Data: retfiles})
	return
}

func GetFsFileInfo(ih metainfo.Hash, fp string) (ret []byte) {
	fp = filepath.ToSlash(fp)
	ret, _ = json.Marshal(DataMsg{Type: "fsfileinfo", Data: nil})
	ihs := ih.HexString()
	if ihs == "" {
		Warn.Println("empty infohash")
		return
	}
	dirpath := filepath.Join(Dirconfig.TrntDir, ihs, fp)

	if !strings.HasPrefix(dirpath, filepath.Join(Dirconfig.TrntDir, ihs)) {
		return
	}
	var retfile FsFileInfo
	inf, err := os.Stat(filepath.FromSlash(dirpath))
	if err != nil {
		Warn.Println("GetFsFileInfo Error", err)
		return
	}
	retfile.Name = inf.Name()
	retfile.IsDir = inf.IsDir()
	retfile.Size = inf.Size()
	retfile.Path = fp
	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "fsfileinfo", Data: retfile})
	return
}

func GetTorrentFiles(ih metainfo.Hash) (ret []byte) {
	ret, _ = json.Marshal(DataMsg{Type: "torrentfiles", Data: nil})
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		return
	}
	if t == nil || t.Info() == nil {
		return
	}

	var retfiles []FileInfo
	for _, file := range t.Files() {
		if file == nil {
			continue
		}
		var retfile FileInfo
		retfile.BytesCompleted = file.BytesCompleted()
		retfile.DisplayPath = file.DisplayPath()
		retfile.Length = file.Length()
		retfile.Offset = file.Offset()
		retfile.Path = file.Path()
		retfile.Priority = byte(file.Priority())

		retfiles = append(retfiles, retfile)
	}

	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "torrentfiles", Data: retfiles})
	return
}

func GetTorrentFileInfo(ih metainfo.Hash, fp string) (ret []byte) {
	ret, _ = json.Marshal(DataMsg{Type: "torrentfileinfo", Data: nil})
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		return
	}
	if t == nil || t.Info() == nil {
		return
	}

	// Get File from Given Torrent
	var file *torrent.File
	for _, f := range t.Files() {
		if f.Path() == fp {
			file = f
			break
		}
	}

	if file == nil {
		return
	}
	var retfile FileInfo
	retfile.BytesCompleted = file.BytesCompleted()
	retfile.DisplayPath = file.DisplayPath()
	retfile.Length = file.Length()
	retfile.Offset = file.Offset()
	retfile.Path = file.Path()
	retfile.Priority = byte(file.Priority())

	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "torrentfileinfo", Data: retfile})
	return
}

func GetTorrentPieceStateRuns(ih metainfo.Hash) (ret []byte) {
	ret, _ = json.Marshal(DataMsg{Type: "torrentpiecestateruns", Data: nil})
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		return
	}
	if t == nil || t.Info() == nil {
		return
	}

	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "torrentpiecestateruns", Data: t.PieceStateRuns()})
	return
}

func GetTorrentKnownSwarm(ih metainfo.Hash) (ret []byte) {
	ret, _ = json.Marshal(DataMsg{Type: "torrentknownswarm", Data: nil})
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		return
	}
	if t == nil || t.Info() == nil {
		return
	}

	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "torrentknownswarm", Data: t.KnownSwarm()})
	return
}

func GetTorrentNumpieces(ih metainfo.Hash) (ret []byte) {
	ret, _ = json.Marshal(DataMsg{Type: "torrentnumpieces", Data: nil})
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		return
	}
	if t == nil || t.Info() == nil {
		return
	}

	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "torrentnumpieces", Data: t.NumPieces()})
	return
}

func GetTorrentMetainfo(ih metainfo.Hash) (ret []byte) {
	ret, _ = json.Marshal(DataMsg{Type: "torrentmetainfo", Data: nil})
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		return
	}
	if t == nil || t.Info() == nil {
		return
	}
	mi := t.Metainfo()
	mi.CreatedBy = "exatorrent"
	var tmi bytes.Buffer
	err := mi.Write(&tmi)
	if err != nil {
		return
	}
	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "torrentmetainfo", Data: tmi.Bytes()})
	return
}

func GetTorrentPeerConns(ih metainfo.Hash) (ret []byte) {
	ret, _ = json.Marshal(DataMsg{Type: "torrentpeerconns", Data: nil})
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		return
	}
	if t == nil || t.Info() == nil {
		return
	}

	var retpeers []PeerInfo
	for _, peerconn := range t.PeerConns() {
		var peerinfo PeerInfo
		peerinfo.RemoteAddr = peerconn.Peer.RemoteAddr.String()
		peerinfo.PeerClientName = fmt.Sprint(peerconn.PeerClientName.Load())
		peerinfo.DownloadRate = peerconn.DownloadRate()
		peerinfo.PeerPreferEncryption = peerconn.PeerPrefersEncryption
		retpeers = append(retpeers, peerinfo)
	}

	ret, _ = json.Marshal(DataMsg{Infohash: ih.HexString(), Type: "torrentpeerconns", Data: retpeers})
	return
}
