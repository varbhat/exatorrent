package db

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/uptrace/bun/extra/bundebug"
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

type ITorrentRepo interface {
	Open(string)
	Close() error
	Exists(metainfo.Hash) bool
	Add(metainfo.Hash) error
	Delete(metainfo.Hash) error
	Start(metainfo.Hash) error
	SetStarted(metainfo.Hash, bool) error
	HasStarted(string) bool
	GetTorrent(metainfo.Hash) (*Torrent, error)
	GetTorrents() ([]*Torrent, error)
}

type ITrackerRepo interface {
	Open(string)
	Close() error
	Add(string)
	Delete(string)
	DeleteN(int)
	DeleteAll()
	Count() int64
	Get() []string
}

type IFileStateRepo interface {
	Open(string)
	Close() error
	Add(string, metainfo.Hash) error
	Get(metainfo.Hash) []string
	Deletefile(string, metainfo.Hash) error
	Delete(metainfo.Hash) error
}

type ILockStateRepo interface {
	Open(string)
	Close() error
	Lock(metainfo.Hash) error
	Unlock(metainfo.Hash) error
	IsLocked(string) bool
}

type IUserRepo interface {
	Open(string)
	Close() error
	Add(string, string, int) error // Username , Password , Usertype
	ChangeType(string, string) error
	Delete(string) error
	UpdatePw(string, string) error
	GetUsers() []*User
	Validate(string, string) (int, bool)
	ValidateToken(string) (string, int, error)
	SetToken(string, string) error
	CheckUserExists(string) bool
}

type ITorrentUserRepo interface {
	Open(string)
	Close() error
	Add(string, metainfo.Hash) error
	Remove(string, metainfo.Hash) error
	RemoveAll(string) error
	RemoveAllMi(metainfo.Hash) error
	HasUser(string, string) bool
	ListTorrents(string) []metainfo.Hash
	ListUsers(metainfo.Hash) []string
}

type IPcRepo interface {
	storage.PieceCompletion
	Delete(metainfo.Hash)
	Open()
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

func InitDb(dbType DBType, dsn string) (db *bun.DB, err error) {
	var sqldb *sql.DB
	switch dbType {
	case Postgres:
		// insecure by default and overwrite by dsn
		sqldb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithInsecure(true), pgdriver.WithDSN(dsn)))
		db = bun.NewDB(sqldb, pgdialect.New())
	case Sqlite:
		sqldb, err = sql.Open(sqliteshim.ShimName, dsn)
		if err != nil {
			return
		}
		db = bun.NewDB(sqldb, sqlitedialect.New())
		db.SetMaxOpenConns(1) // solve sqlite database is locked error also no need to use a lock
	default:
		err = fmt.Errorf("unsupported database type: %s", dbType)
		return
	}

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.FromEnv(),
	))

	err = db.Ping()
	if err != nil {
		return
	}
	return
}

//go:generate go run github.com/dmarkham/enumer -type=DBType
type DBType int

const (
	Postgres DBType = iota
	Sqlite
)
