package core

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/disk"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 31457280 // 30 MB
)

func SocketAPI(w http.ResponseWriter, r *http.Request) {
	username, usertype, _, err := authHelper(w, r)
	if err != nil {
		Warn.Printf("%s (%s)\n", err, r.RemoteAddr)
		return
	}
	var admin bool
	if usertype == 1 {
		admin = true
	} else if usertype == -1 {
		Err.Println("Disabled User denied to Connect", username)
		http.Error(w, "User Disabled", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Warn.Println(err)
		return
	}
	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { _ = conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	Uc := NewUserConn(username, conn, admin)

	// Ping Handler
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
		}()

		var perr error

		for {
			<-ticker.C
			Uc.Sendmu.Lock()
			_ = Uc.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			perr = Uc.Conn.WriteMessage(websocket.PingMessage, nil)
			Uc.Sendmu.Unlock()
			if perr != nil {
				Uc.Close()
				return
			}
		}
	}()

	var Req []byte
	for {
		_, Req, err = Uc.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				Warn.Printf("WebSocket Unexpected Close: %v", err)
			}
			Uc.Close()
			return
		}
		var resp ConReq
		err = json.Unmarshal(Req, &resp)
		if err != nil {
			_ = Uc.SendMsg("resp", "error", "incorrect request")
			continue
		}
		go wshandler(Uc, &resp)
	}
}

func wshandler(uc *UserConn, req *ConReq) {
	if len(req.Command) == 0 {
		_ = uc.SendMsg("resp", "error", "no command sent")
		return
	}

	// Admin Operations
	if uc.IsAdmin && req.Aop == 1 {
		switch req.Command {
		// parse these in normal block
		case "abandontorrent":
		case "removetorrent":
		case "addtrackerstotorrent":
		case "starttorrent":
		case "stoptorrent":
		case "startfile":
		case "stopfile":
		case "deletefilepath":
		case "deletetorrent":
		//
		case "adduser":
			if !(len(req.Data1) > 5 || len(req.Data2) > 5) {
				_ = uc.SendMsg("resp", "error", "length of username and password must be more than 5")
				return
			}
			var err error
			if req.Data3 == "admin" {
				err = Engine.UDb.Add(req.Data1, req.Data2, 1) // Admin
			} else if req.Data3 == "user" {
				err = Engine.UDb.Add(req.Data1, req.Data2, 0) // User
			} else if req.Data3 == "disabled" {
				err = Engine.UDb.Add(req.Data1, req.Data2, -1) // Disabled
			} else {
				_ = uc.SendMsg("resp", "error", "incorrect adduser request")
				return
			}
			if err != nil {
				_ = uc.SendMsg("resp", "error", "Error Adding User to Database")
				return
			}
			_ = uc.SendMsg("resp", "success", req.Data1+" User Added")
			Info.Println("New User ", req.Data1, " added by ", uc.Username)
			return
		case "removeuser":
			if !(len(req.Data1) > 5) {
				_ = uc.SendMsg("resp", "error", "length of username  must be more than 5")
				return
			}
			if req.Data1 == uc.Username {
				_ = uc.SendMsg("resp", "error", "you cannot remove yourself")
				return
			}
			err := Engine.UDb.Delete(req.Data1)
			if err != nil {
				_ = uc.SendMsg("resp", "error", "Error Removing User from Database")
				return
			}
			err = Engine.TUDb.RemoveAll(req.Data1)
			if err != nil {
				_ = uc.SendMsg("resp", "error", "Error Deleting User Records from Database")
				return
			}
			MainHub.RemoveUser(req.Data1)
			_ = uc.SendMsg("resp", "success", req.Data1+" User Removed")
			Info.Println("User ", req.Data1, " is removed by ", uc.Username)
			return
		case "getalltorrents":
			uc.Streamers.Inc()
			uc.Stream.Lock()
			defer uc.Streamers.Dec()
			defer uc.Stream.Unlock()

			files, ferr := os.ReadDir(Dirconfig.TrntDir)
			if ferr != nil {
				Warn.Println(ferr)
				return
			}

			var lt []metainfo.Hash = make([]metainfo.Hash, 0)
			for _, file := range files {
				if file.IsDir() {
					tm, terr := MetafromHex(file.Name())
					if terr != nil {
						Warn.Println(terr)
						Warn.Println("Non Torrent Directories found in ", Dirconfig.TrntDir, file, file.Name())
					}
					lt = append(lt, tm)
				} else {
					Warn.Println("Non Torrent Directories found in ", Dirconfig.TrntDir, file, file.Name())
				}
			}

			var err error
			Info.Println("Starting getalltorrents for ", uc.Username)
			for uc.Streamers.Get() == 1 {
				err = uc.Send(GetTorrents(lt))
				if err != nil {
					return
				}
				if uc.Streamers.Get() == 1 {
					time.Sleep(time.Second * 5) // Stream Every 5 Seconds
				}
			}
			Info.Println("Stopped getalltorrents for ", uc.Username)
			return
		case "listalltorrents":
			files, ferr := os.ReadDir(Dirconfig.TrntDir)
			if ferr != nil {
				Warn.Println(ferr)
				return
			}
			var lt []metainfo.Hash = make([]metainfo.Hash, 0)
			for _, file := range files {
				if file.IsDir() {
					tm, terr := MetafromHex(file.Name())
					if terr != nil {
						Warn.Println(terr)
						Warn.Println("Non Torrent Directories found in ", Dirconfig.TrntDir, file, file.Name())
					}
					lt = append(lt, tm)
				} else {
					Warn.Println("Non Torrent Directories found in ", Dirconfig.TrntDir, file, file.Name())
				}
			}
			_ = uc.Send(GetTorrents(lt))
			return
		case "listtorrentsforuser":
			if req.Data1 != "" {
				ret, _ := json.Marshal(DataMsg{Type: "torrentsforuser", Data: Engine.TUDb.ListTorrents(req.Data1)})
				_ = uc.Send(ret)
			} else {
				_ = uc.SendMsg("resp", "error", "username is empty")
			}
			return
		case "getusers":
			retusers := Engine.UDb.GetUsers()
			ret, _ := json.Marshal(DataMsg{Type: "users", Data: retusers})
			_ = uc.Send(ret)
			return
		case "updatepw":
			defer func() {
				if r := recover(); r != nil { // uuid may panic
					_ = uc.SendMsg("resp", "error", "uuid error")
					Warn.Println("uuid error")
					return
				}
			}()
			if req.Data1 == "" {
				_ = uc.SendMsg("resp", "error", "request error")
				return
			}
			err := Engine.UDb.UpdatePw(req.Data1, req.Data2)
			if err != nil {
				_ = uc.SendMsg("resp", "error", err.Error())
			}
			MainHub.RemoveUser(req.Data1)
			err = Engine.UDb.SetToken(req.Data1, uuid.New().String())
			if err != nil {
				_ = uc.SendMsg("resp", "error", err.Error())
			}
			_ = uc.SendMsg("resp", "success", "Password of "+req.Data1+" updated")
			return
		case "revoketoken":
			defer func() {
				if r := recover(); r != nil { // uuid may panic
					_ = uc.SendMsg("resp", "error", "uuid error")
					Warn.Println("uuid error")
					return
				}
			}()
			if req.Data1 == "" {
				_ = uc.SendMsg("resp", "error", "request error")
				return
			}
			MainHub.RemoveUser(req.Data1)
			err := Engine.UDb.SetToken(req.Data1, uuid.New().String())
			if err != nil {
				_ = uc.SendMsg("resp", "error", err.Error())
			}
			_ = uc.SendMsg("resp", "success", "revoked token of "+req.Data1)
			return
		case "changeusertype":
			Info.Println(uc.Username, "has requested to change usertype of ", req.Data1, " to ", req.Data2)
			if req.Data1 == "" {
				_ = uc.SendMsg("resp", "error", "username can be empty")
				Warn.Println("Usertype change request of ", uc.Username, " failed")
				return
			}
			if req.Data1 == uc.Username {
				_ = uc.SendMsg("resp", "error", "you can't change usertype of yourself")
				Warn.Println("Usertype change request of ", uc.Username, " failed")
				return
			}
			MainHub.RemoveUser(req.Data1)
			err := Engine.UDb.ChangeType(req.Data1, req.Data2)
			if err != nil {
				_ = uc.SendMsg("resp", "error", err.Error())
			}
			_ = uc.SendMsg("resp", "success", "changed usertype of "+req.Data1+" to "+req.Data2)
			Info.Println(uc.Username, "has changed usertype of ", req.Data1, " to ", req.Data2)
			return
		case "machinfo":
			ret, _ := json.Marshal(DataMsg{Type: "machinfo", Data: MachInfo})
			_ = uc.Send(ret)
			return
		case "machstats":
			MachStats.LoadStats(".")
			ret, _ := json.Marshal(DataMsg{Type: "machstats", Data: MachStats})
			_ = uc.Send(ret)
			return
		case "torcstatus":
			var torcstatus bytes.Buffer
			Engine.Torc.WriteStatus(&torcstatus)
			ret, _ := json.Marshal(DataMsg{Type: "torcstatus", Data: torcstatus.String()})
			_ = uc.Send(ret)
			return
		case "getconfig":
			Configmu.Lock()
			ret, _ := json.Marshal(DataMsg{Type: "engconf", Data: Engine.Econfig})
			Configmu.Unlock()
			_ = uc.Send(ret)
			return
		case "updateconfig":
			if req.Data1 == "" {
				_ = uc.SendMsg("resp", "error", "empty config file")
				return
			}
			configfile, berr := base64.StdEncoding.DecodeString(req.Data1)
			if berr != nil {
				_ = uc.SendMsg("resp", "error", "error decoding config file")
				return
			}
			var newconfig EngConfig
			berr = json.Unmarshal(configfile, &newconfig)
			if berr != nil {
				_ = uc.SendMsg("resp", "error", "error decoding config file")
				return
			}
			Configmu.Lock()
			if Engine.Econfig.ListenCompletion != newconfig.ListenCompletion {
				if !newconfig.ListenCompletion {
					for _, eachchan := range Engine.onCloseMap {
						if eachchan != nil {
							eachchan.Set()
							eachchan.Clear()
						}
					}
				} else {
					trntlist := Engine.Torc.Torrents()
					for _, eachtrnt := range trntlist {
						if eachtrnt != nil {
							infohash := eachtrnt.InfoHash()
							go func(eachtrnt *torrent.Torrent) {
								if _, ok := Engine.onCloseMap[infohash]; !ok {
									Engine.onCloseMap[infohash] = &eachtrnt.Complete
									Info.Println("Listening for Completion of Torrent ", infohash)
									<-eachtrnt.Complete.On()
									delete(Engine.onCloseMap, infohash)

									_, err := Engine.TorDb.GetTorrent(infohash)
									if err != nil {
										Info.Println(infohash, " Removed")
									} else {
										Info.Println(infohash, " Completed")
										hpu := Engine.Econfig.GetHPU()
										if hpu != "" {
											trntname := ""
											if eachtrnt != nil {
												trntname = eachtrnt.Name()
											}
											sendPostReq(infohash, hpu, trntname)
										}
									}
								}
							}(eachtrnt)
						}
					}
				}
			}
			Engine.Econfig = newconfig
			_ = Engine.Econfig.WriteConfig()
			Info.Println("Torrent Configuration has been Updated by ", uc.Username)
			Configmu.Unlock()
			_ = uc.SendMsg("resp", "success", "New Torrent Config File has been set successfully")
			return
		case "listusersfortorrent":
			ih, err := MetafromHex(req.Data1)
			if err != nil {
				_ = uc.SendMsg("resp", "error", "listusersfortorrent: infohash couldn't be parsed "+req.Data1)
				return
			}
			ret, _ := json.Marshal(DataMsg{Type: "usersfortorrent", Infohash: ih.HexString(), Data: Engine.TUDb.ListUsers(ih)})
			_ = uc.Send(ret)
			return
		case "listuserconns":
			_ = uc.Send(MainHub.ListUsers())
			return
		case "kickuser":
			MainHub.RemoveUser(req.Data1)
			_ = uc.SendMsg("resp", "success", "kicked "+req.Data1)
			return
		case "changedataload":
			ih, err := MetafromHex(req.Data1)
			if err != nil {
				_ = uc.SendMsg("resp", "error", "stopfile: infohash couldn't be parsed "+req.Data1)
				return
			}
			t, ok := Engine.Torc.Torrent(ih)
			if !ok {
				_ = uc.SendMsgU("nfn", ih.HexString(), "error", "Torrent Not Present!")
				Warn.Println("Error fetching torrent of infohash ", ih, " from the client")
				return
			}

			if t.Info() != nil {
				if req.Data2 == "upload" {
					if req.Data3 == "disallow" {
						t.DisallowDataUpload()
						_ = uc.SendMsg("resp", "success", "Disallowed Data Upload for torrent "+ih.HexString())
						return
					} else if req.Data3 == "allow" {
						t.AllowDataUpload()
						_ = uc.SendMsg("resp", "success", "Allowed Data Upload for torrent "+ih.HexString())
						return
					} else {
						_ = uc.SendMsg("resp", "error", "invalid request")
						return
					}
				} else if req.Data2 == "download" {
					if req.Data3 == "disallow" {
						t.DisallowDataDownload()
						_ = uc.SendMsg("resp", "success", "Disallowed Data Download for torrent "+ih.HexString())
						return
					} else if req.Data3 == "allow" {
						t.AllowDataDownload()
						_ = uc.SendMsg("resp", "success", "Allowed Data Download for torrent "+ih.HexString())
						return
					} else {
						_ = uc.SendMsg("resp", "error", "invalid request")
						return
					}
				} else {
					_ = uc.SendMsg("resp", "error", "invalid request")
					return
				}
			}
			return
		case "nooftrackersintrackerdb":
			ret, _ := json.Marshal(DataMsg{Type: "nooftrackersintrackerdb", Data: Engine.TrackerDB.Count()})
			_ = uc.Send(ret)
			return
		case "deletetrackersintrackerdb":
			if req.Data1 == "all" {
				Engine.TrackerDB.DeleteAll()
				Info.Println("Deleted All Trackers from TrackerDB")
				_ = uc.SendMsg("resp", "success", "Deleted All Trackers from TrackerDB")
				return
			}
			notobedeleted, err := strconv.Atoi(req.Data1)
			if err != nil {
				_ = uc.SendMsg("resp", "error", "invalid request")
				return
			}
			Engine.TrackerDB.DeleteN(notobedeleted)
			Info.Println("Deleted ", req.Data1, " no of Trackers from TrackerDB")
			_ = uc.SendMsg("resp", "success", "Deleted "+req.Data1+" trackers from TrackerDB")
			return
		case "trackerdbrefresh":
			Info.Println("TrackerDB Refresh command has been issued by", uc.Username)
			updatetrackers()
			_ = uc.SendMsg("resp", "success", "Fetched Trackers from Tracker URLs and Updated TrackerDB")
			return
		case "stoponseedratio":
			Info.Println("stoponseedratio command has been issued by ", uc.Username)
			if req.Data1 == "" {
				stoponseedratio(Engine.Econfig.GetGSR())
				_ = uc.SendMsg("resp", "success", "Stopped Torrents which have reached  seedratio defined in config")
				return

			}
			if sr, err := strconv.ParseFloat(req.Data1, 64); err == nil {
				stoponseedratio(sr)
				_ = uc.SendMsg("resp", "success", "Stopped Torrents which have reached  seedratio = "+req.Data1)
				return
			}
			_ = uc.SendMsg("resp", "error", "Invalid request")
			return
		default:
			_ = uc.SendMsg("resp", "error", "invalid admin command")
			return
		}
	}

	switch req.Command {
	case "addmagnet":
		tspec, terr := torrent.TorrentSpecFromMagnetUri(req.Data1)
		if terr != nil {
			_ = uc.SendMsg("resp", "error", "incorrect torrent spec")
			return
		}
		if req.Data2 == "true" {
			go AddfromSpec(uc.Username, tspec, true, false)
		} else {
			go AddfromSpec(uc.Username, tspec, false, false)
		}
		_ = uc.SendMsgU("resp", "success", tspec.InfoHash.HexString(), "torrent spec added")
		return
	case "addtorrent":
		tspec, terr := SpecfromB64String(req.Data1)
		if terr != nil {
			_ = uc.SendMsg("resp", "error", "incorrect torrent spec")
			return
		}
		if req.Data2 == "true" {
			go AddfromSpec(uc.Username, tspec, true, false)
		} else {
			go AddfromSpec(uc.Username, tspec, false, false)
		}
		_ = uc.SendMsgU("resp", "success", tspec.InfoHash.HexString(), "torrent spec added")
		return
	case "addinfohash":
		ih, terr := MetafromHex(req.Data1)
		if terr != nil {
			_ = uc.SendMsg("resp", "error", "incorrect Infohash")
			return
		}
		tspec := &torrent.TorrentSpec{InfoHash: ih}
		if req.Data2 == "true" {
			go AddfromSpec(uc.Username, tspec, true, false)
		} else {
			go AddfromSpec(uc.Username, tspec, false, false)
		}
		_ = uc.SendMsgU("resp", "success", ih.HexString(), "torrent spec added")
		return
	case "starttorrent":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "starttorrent: infohash couldn't be parsed "+req.Data1)
			return
		}
		if uc.IsAdmin && req.Aop == 1 {
			go StartTorrent("", ih, false)
			return
		}
		go StartTorrent(uc.Username, ih, false)
		return
	case "stoptorrent":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "stoptorrent: infohash couldn't be parsed "+req.Data1)
			return
		}
		if uc.IsAdmin && req.Aop == 1 {
			go StopTorrent("", ih)
			return
		}
		go StopTorrent(uc.Username, ih)
		return
	case "startfile":
		if req.Data2 == "" {
			_ = uc.SendMsg("resp", "error", "no path provided")
			return
		}
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "startfile: infohash couldn't be parsed "+req.Data1)
			return
		}
		if uc.IsAdmin && req.Aop == 1 {
			go StartFile("", ih, req.Data2)
			return
		}
		go StartFile(uc.Username, ih, req.Data2)
		return
	case "stopfile":
		if req.Data2 == "" {
			_ = uc.SendMsg("resp", "error", "no path provided")
			return
		}
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "stopfile: infohash couldn't be parsed "+req.Data1)
			return
		}
		if uc.IsAdmin && req.Aop == 1 {
			go StopFile("", ih, req.Data2)
			return
		}
		go StopFile(uc.Username, ih, req.Data2)
		return
	case "deletefilepath":
		if req.Data2 == "" {
			_ = uc.SendMsg("resp", "error", "no path provided")
			return
		}
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "deletefilepath: infohash couldn't be parsed "+req.Data1)
			return
		}
		if uc.IsAdmin && req.Aop == 1 {
			go DeleteFilePath("", ih, req.Data2)
			return
		}
		go DeleteFilePath(uc.Username, ih, req.Data2)
		return
	case "deletetorrent":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "deletetorrent: infohash couldn't be parsed "+req.Data1)
			return
		}
		if uc.IsAdmin && req.Aop == 1 {
			go DeleteTorrent("", ih)
			return
		}
		go DeleteTorrent(uc.Username, ih)
		uc.StopStream()
		return
	case "removetorrent":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "removetorrent: infohash couldn't be parsed "+req.Data1)
			return
		}
		if uc.IsAdmin && req.Aop == 1 {
			go RemoveTorrent("", ih)
			return
		}
		go RemoveTorrent(uc.Username, ih)
		return
	case "abandontorrent":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "abandontorrent: infohash couldn't be parsed "+req.Data1)
			return
		}
		if uc.IsAdmin && req.Aop == 1 {
			go AbandonTorrent(req.Data2, ih)
			return
		}
		go AbandonTorrent(uc.Username, ih)
		return
	case "gettorrents":
		uc.Streamers.Inc()
		uc.Stream.Lock()
		defer uc.Streamers.Dec()
		defer uc.Stream.Unlock()

		lt := Engine.TUDb.ListTorrents(uc.Username)
		var err error
		Info.Println("Starting gettorrents for ", uc.Username)
		for uc.Streamers.Get() == 1 {
			err = uc.Send(GetTorrents(lt))
			if err != nil {
				return
			}
			if uc.Streamers.Get() == 1 {
				time.Sleep(time.Second * 5) // Stream Every 5 Seconds
			}
		}
		Info.Println("Stopped gettorrents for ", uc.Username)
		return
	case "listtorrents":
		_ = uc.Send(GetTorrents(Engine.TUDb.ListTorrents(uc.Username)))
		return
	case "gettorrentinfo":
		uc.Streamers.Inc()
		uc.Stream.Lock()
		defer uc.Streamers.Dec()
		defer uc.Stream.Unlock()

		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentinfo: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		Info.Println("Starting gettorrentinfo for ", uc.Username)

		for uc.Streamers.Get() == 1 {
			err = uc.Send(GetTorrentInfo(ih))
			if err != nil {
				return
			}
			if uc.Streamers.Get() == 1 {
				time.Sleep(time.Second * 5) // Stream Every 5 Seconds
			}
		}
		Info.Println("Stopped gettorrentinfo for ", uc.Username)
		return
	case "listtorrentinfo":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "listtorrentinfo: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		_ = uc.Send(GetTorrentInfo(ih))
		return
	case "gettorrentstats":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentstats: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentStats(ih))
		if err != nil {
			return
		}
		return
	case "gettorrentinfostat":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentinfostat: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentInfoStat(ih))
		if err != nil {
			return
		}
		return
	case "gettorrentfiles":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentfiles: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentFiles(ih))
		if err != nil {
			return
		}
		return
	case "getfsdirinfo":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "getfsdirinfo: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetFsDirInfo(ih, filepath.Clean(req.Data2)))
		if err != nil {
			return
		}
		return
	case "getfsfileinfo":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "getfsfileinfo: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetFsFileInfo(ih, filepath.Clean(req.Data2)))
		if err != nil {
			return
		}
		return
	case "gettorrentfileinfo":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentfileinfo: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentFileInfo(ih, filepath.Clean(req.Data2)))
		if err != nil {
			return
		}
		return
	case "gettorrentpiecestateruns":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentpiecestateruns: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentPieceStateRuns(ih))
		if err != nil {
			return
		}
		return
	case "istorrentlocked":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "istorrentlocked: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}

		ret, _ := json.Marshal(DataMsg{Type: "torrentlockstate", Infohash: ih.HexString(), Data: Engine.LsDb.IsLocked(ih.HexString())})
		_ = uc.Send(ret)
		return
	case "toggletorrentlock":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "toggletorrentlock: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		} else {
			_, err = os.Stat(filepath.Join(Dirconfig.TrntDir, ih.HexString()))
			if err != nil {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}

		if Engine.LsDb.IsLocked(ih.HexString()) {
			_ = Engine.LsDb.Unlock(ih)
			Info.Println("lock of torrent ", ih.HexString(), " is set to false by ", uc.Username)
			_ = uc.SendMsgU("resp", "success", ih.HexString(), "lock of torrent "+ih.HexString()+" is set to false")
		} else {
			_ = Engine.LsDb.Lock(ih)
			Info.Println("lock of torrent ", ih.HexString(), " is set to true by ", uc.Username)
			_ = uc.SendMsgU("resp", "success", ih.HexString(), "lock of torrent "+ih.HexString()+" is set to true")
		}
		return
	case "addtrackerstotorrent":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "addtrackerstotorrent: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		if req.Data2 == "" {
			_ = uc.SendMsg("resp", "error", "Invalid Request")
		}
		// Add Trackers to txtlines string slice
		trackerlist, berr := base64.StdEncoding.DecodeString(req.Data2)
		if berr != nil {
			_ = uc.SendMsg("resp", "error", "trackerlist error")
			return
		}
		tlr := bytes.NewReader(trackerlist)
		scanner := bufio.NewScanner(tlr)
		scanner.Split(bufio.ScanLines)

		var trackerno int
		var al []string
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			al = append(al, line)
			trackerno++
		}
		Info.Println("Read ", trackerno, " Trackers from Trackerlist uploaded by ", uc.Username)
		if uc.IsAdmin && req.Aop == 1 {
			go AddTrackerstoTorrent("", ih, [][]string{al})
		} else {
			go AddTrackerstoTorrent(uc.Username, ih, [][]string{al})
		}
		return
	case "updatepw":
		err := Engine.UDb.UpdatePw(uc.Username, req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", err.Error())
		}
		return
	case "stopstream":
		uc.StopStream()
		return
	case "diskusage":
		var diskusage DiskUsageStat
		if stat, err := disk.Usage(Dirconfig.TrntDir); err == nil {
			diskusage.UsedPercent = stat.UsedPercent
			diskusage.Free = stat.Free
			diskusage.Total = stat.Total
			diskusage.Used = stat.Used
		}
		ret, _ := json.Marshal(DataMsg{Type: "diskusage", Data: diskusage})
		_ = uc.Send(ret)
		return
	case "gettorrentknownswarm":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentknownswarm: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentKnownSwarm(ih))
		if err != nil {
			return
		}
		return
	case "gettorrentnumpieces":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentnumpieces: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentNumpieces(ih))
		if err != nil {
			return
		}
		return
	case "gettorrentmetainfo":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentmetainfo: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentMetainfo(ih))
		if err != nil {
			return
		}
		return
	case "gettorrentpeerconns":
		ih, err := MetafromHex(req.Data1)
		if err != nil {
			_ = uc.SendMsg("resp", "error", "gettorrentpeerconns: infohash couldn't be parsed "+req.Data1)
			return
		}

		if !uc.IsAdmin {
			if !Engine.TUDb.HasUser(uc.Username, ih.HexString()) {
				_ = uc.SendMsg("resp", "error", "torrent doesn't exist")
				return
			}
		}
		err = uc.Send(GetTorrentPeerConns(ih))
		if err != nil {
			return
		}
		return
	case "version":
		ret, _ := json.Marshal(DataMsg{Type: "version", Data: Version})
		_ = uc.Send(ret)
		return
	default:
		_ = uc.SendMsg("resp", "error", "invalid command")
	}
}
