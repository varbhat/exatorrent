package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

var DbL = log.New(os.Stderr, "[DB]  ", log.LstdFlags) // Database Logger

// MetafromHex returns metainfo.Hash from given infohash string
func MetafromHex(infohash string) (h metainfo.Hash, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error parsing string to InfoHash")
		}
	}()

	h = metainfo.NewHashFromHex(infohash)

	return h, nil
}

// Interfaces

type TorrentDb interface {
	Open(string)
	Close()
	Exists(metainfo.Hash) bool
	Add(metainfo.Hash) error
	Delete(metainfo.Hash) error
	Start(metainfo.Hash) error
	SetStarted(metainfo.Hash, bool) error
	HasStarted(string) bool
	GetTorrent(metainfo.Hash) (*Torrent, error)
	GetTorrents() ([]*Torrent, error)
}

type TrackerDb interface {
	Open(string)
	Close()
	Add(string)
	Delete(string)
	DeleteN(int)
	DeleteAll()
	Count() int64
	Get() []string
}

type FileStateDb interface {
	Open(string)
	Close()
	Add(string, metainfo.Hash) error
	Get(metainfo.Hash) []string
	Deletefile(string, metainfo.Hash) error
	Delete(metainfo.Hash) error
}

type LockStateDb interface {
	Open(string)
	Close()
	Lock(metainfo.Hash) error
	Unlock(metainfo.Hash) error
	IsLocked(string) bool
}

type UserDb interface {
	Open(string)
	Close()
	Add(string, string, int) error // Username , Password , Usertype
	ChangeType(string, string) error
	Delete(string) error
	UpdatePw(string, string) error
	GetUsers() []*User
	Validate(string, string) (int, bool)
	ValidateToken(string) (string, int, error)
	SetToken(string, string) error
}

type TorrentUserDb interface {
	Open(string)
	Close()
	Add(string, metainfo.Hash) error
	Remove(string, metainfo.Hash) error
	RemoveAll(string) error
	RemoveAllMi(metainfo.Hash) error
	HasUser(string, string) bool
	ListTorrents(string) []metainfo.Hash
	ListUsers(metainfo.Hash) []string
}

type PcDb interface {
	storage.PieceCompletion
	Delete(metainfo.Hash)
}

// Struct

type Torrent struct {
	Infohash  metainfo.Hash
	Started   bool
	AddedAt   time.Time
	StartedAt time.Time
}

type User struct {
	Username  string
	Password  string `json:"-"`
	Token     string
	UserType  int // 0 for User,1 for Admin,-1 for Disabled
	CreatedAt time.Time
}
