//go:build cgo
// +build cgo

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/anacrolix/torrent/metainfo"
)

type TorrentEntity struct {
	InfoHash  string
	Started   bool
	AddedAt   time.Time
	StartedAt time.Time
}

type Scannable interface {
	Scan(...any) error
}

type Sqlite3Db struct {
	Db *sql.DB
	mu sync.Mutex
}

func (db *Sqlite3Db) Open(fp string) {
	var err error

	db.Db, err = sql.Open("sqlite3", fp)
	db.Db.SetMaxOpenConns(1)

	if err != nil {
		DbL.Fatalln(err)
	}

	_, err = db.Db.Exec(`create table if not exists torrent (infohash text primary key,started boolean,addedat text,startedat text);`)

	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *Sqlite3Db) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.Db.Close()
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *Sqlite3Db) Exists(ih metainfo.Hash) (ret bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.Db.QueryRow(`select exists (select 1 from torrent where infohash=?) ;`, ih.HexString())
	err := row.Err()

	if err != nil {
		return false
	}

	err = row.Scan(&ret)
	if err != nil {
		return false
	}
	return ret
}

func (db *Sqlite3Db) IsLocked(ih string) (ret bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	row := db.Db.QueryRow(`select locked from torrent where infohash=?;`, ih)
	err := row.Err()
	if err != nil {
		return false
	}
	err = row.Scan(&ret)
	if err != nil {
		return false
	}
	return ret
}

func (db *Sqlite3Db) HasStarted(ih string) (ret bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	row := db.Db.QueryRow(`select started from torrent where infohash=?;`, ih)
	err := row.Err()
	if err != nil {
		return false
	}
	err = row.Scan(&ret)
	if err != nil {
		return false
	}
	return ret
}

func (db *Sqlite3Db) SetLocked(ih string, b bool) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`update torrent set locked=? where infohash=?;`, b, ih)
	return
}

func (db *Sqlite3Db) Add(ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	tn := time.Now().Format(time.RFC3339)
	_, err = db.Db.Exec(`insert into torrent (infohash,started,addedat,startedat) values (?,?,?,?) on conflict (infohash) do update set startedat=?;`, ih.HexString(), 0, tn, tn, tn)
	return
}

func (db *Sqlite3Db) Delete(ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`delete from torrent where infohash=?;`, ih.HexString())
	return
}

func (db *Sqlite3Db) Start(ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	tn := time.Now().Format(time.RFC3339)
	_, err = db.Db.Exec(`update torrent set started=?,startedat=? where infohash=?;`, 1, tn, ih.HexString())
	return
}

func (db *Sqlite3Db) SetStarted(ih metainfo.Hash, inp bool) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`update torrent set started=? where infohash=?;`, inp, ih.HexString())
	return
}

func (db *Sqlite3Db) GetTorrent(ih metainfo.Hash) (*Torrent, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	row := db.Db.QueryRow(`select * from torrent where infohash=?;`, ih.HexString())
	err := row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		err = fmt.Errorf("Torrent doesn't exist")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	entity, err := parseRow(row)
	if err != nil {
		return nil, err
	}
	return &Torrent{
		Infohash:  ih,
		Started:   entity.Started,
		AddedAt:   entity.AddedAt,
		StartedAt: entity.StartedAt,
	}, nil
}

func (db *Sqlite3Db) GetTorrents() (torrentList []*Torrent, err error) {
	torrentList = make([]*Torrent, 0)

	db.mu.Lock()
	defer db.mu.Unlock()
	rows, err := db.Db.Query(`select * from torrent;`)
	if err != nil {
		return
	}
	var entity TorrentEntity
	var ih metainfo.Hash
	for rows.Next() {
		entity, err = parseRow(rows)
		if err != nil {
			DbL.Println(err)
		}
		ih, err = MetafromHex(entity.InfoHash)
		if err != nil {
			return
		}
		torrentList = append(torrentList, &Torrent{
			Infohash:  ih,
			Started:   entity.Started,
			AddedAt:   entity.AddedAt,
			StartedAt: entity.StartedAt,
		})
	}

	return torrentList, nil
}

func parseRow(row Scannable) (entity TorrentEntity, err error) {
	var addedAt string
	var startedAt string
	err = row.Scan(&entity.InfoHash, &entity.Started, &addedAt, &startedAt)
	if err != nil {
		return
	}
	entity.AddedAt, err = time.Parse(time.RFC3339, addedAt)
	if err != nil {
		return
	}
	entity.StartedAt, err = time.Parse(time.RFC3339, startedAt)
	if err != nil {
		return
	}
	return
}
