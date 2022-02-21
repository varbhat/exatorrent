package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/anacrolix/torrent/metainfo"

	"github.com/gorilla/websocket"
)

// Hub
type Hub struct {
	sync.RWMutex
	Conns map[string]*UserConn
}

func (h *Hub) Add(Uc *UserConn) {
	conn, ok := h.Conns[Uc.Username]
	if ok && conn != nil {
		conn.Close()
	}
	h.Lock()
	h.Conns[Uc.Username] = Uc
	h.Unlock()

	Info.Printf("User %s (%s) Connected\n", Uc.Username, Uc.Conn.RemoteAddr().String())
}

func (h *Hub) SendMsg(User string, Type string, State string, Resp string) {
	if User != "" {
		conn, ok := h.Conns[User]
		if ok && conn != nil {
			_ = conn.SendMsg(Type, State, Resp)
		}
	}
}

func (h *Hub) SendMsgU(User string, Type string, Infohash string, State string, Resp string) {
	if User != "" {
		conn, ok := h.Conns[User]
		if ok && conn != nil {
			_ = conn.SendMsgU(Type, State, Infohash, Resp)
		}
	}
}

func (h *Hub) Remove(Uc *UserConn) {
	if Uc == nil {
		return
	}
	h.Lock()
	defer h.Unlock()
	conn, ok := h.Conns[Uc.Username]
	if ok && conn != nil {
		if conn.Time == Uc.Time {
			delete(h.Conns, Uc.Username)
			Info.Printf("User %s (%s) Disconnected\n", Uc.Username, Uc.Conn.RemoteAddr().String())
		}
	}
}

func (h *Hub) RemoveUser(Username string) {
	conn, ok := h.Conns[Username]
	if ok && conn != nil {
		conn.Close()
	}
}

func (h *Hub) ListUsers() (ret []byte) {
	var userconnmsg []*UserConnMsg
	h.Lock()
	defer h.Unlock()
	for name, user := range h.Conns {
		var usermsg UserConnMsg
		if user != nil {
			usermsg.Username = name
			usermsg.IsAdmin = user.IsAdmin
			usermsg.Time = user.Time
		}
		userconnmsg = append(userconnmsg, &usermsg)
	}
	ret, _ = json.Marshal(DataMsg{Type: "userconn", Data: userconnmsg})
	return
}

var MainHub Hub = Hub{
	RWMutex: sync.RWMutex{},
	Conns:   make(map[string]*UserConn),
}

// UserConn
type UserConn struct {
	Sendmu    sync.Mutex
	Username  string
	IsAdmin   bool
	Time      time.Time
	Conn      *websocket.Conn
	Stream    sync.Mutex
	Streamers MutInt
}

func NewUserConn(Username string, Conn *websocket.Conn, IsAdmin bool) (uc *UserConn) {
	uc = &UserConn{
		Username: Username,
		Conn:     Conn,
		IsAdmin:  IsAdmin,
		Time:     time.Now(),
	}
	MainHub.Add(uc)
	return
}

func (uc *UserConn) SendMsg(Type string, State string, Msg string) (err error) {
	resp, _ := json.Marshal(Resp{Type: Type, State: State, Msg: Msg})
	err = uc.Send(resp)
	return
}

func (uc *UserConn) SendMsgU(Type string, State string, Infohash string, Msg string) (err error) {
	resp, _ := json.Marshal(Resp{Type: Type, State: State, Infohash: Infohash, Msg: Msg})
	err = uc.Send(resp)
	return
}

func (uc *UserConn) Send(v []byte) (err error) {
	uc.Sendmu.Lock()
	_ = uc.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	err = uc.Conn.WriteMessage(websocket.TextMessage, v)
	uc.Sendmu.Unlock()
	if err != nil {
		Err.Println(err)
		uc.Close()
		return
	}
	return
}

func (uc *UserConn) StopStream() {
	uc.Streamers.Inc()
	uc.Stream.Lock()
	Info.Println("Stopped Stream for ", uc.Username)
	uc.Stream.Unlock()
	uc.Streamers.Dec()
}

func (uc *UserConn) Close() {
	uc.Sendmu.Lock()
	_ = uc.Conn.Close()
	uc.Sendmu.Unlock()
	MainHub.Remove(uc)
}

func sendPostReq(h metainfo.Hash, url string, name string) {
	Info.Println("Torrent ", h, " has completed. Sending POST request to ", url)
	postrequest := struct {
		Metainfo metainfo.Hash `json:"metainfo"`
		Name     string        `json:"name"`
		State    string        `json:"state"`
		Time     time.Time     `json:"time"`
	}{
		Metainfo: h,
		Name:     name,
		Time:     time.Now(),
		State:    "torrent-completed-exatorrent",
	}

	jsonData, err := json.Marshal(postrequest)

	if err != nil {
		Warn.Println(err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		Warn.Println("POST Request failed to Send. Hook failed")
		Warn.Println(err)
		return
	}

	if resp != nil {
		resp.Body.Close()
	}

	Info.Println("POST Request Sent. Hook Succeeded")
}
