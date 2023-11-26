//go:build cgo
// +build cgo

package db

import (
	"fmt"
	"sync"
	"time"

	sqlite "github.com/go-llsqlite/crawshaw"
	"github.com/go-llsqlite/crawshaw/sqlitex"

	"github.com/anacrolix/torrent/metainfo"
)

type Sqlite3Db struct {
	Db *sqlite.Conn
	mu sync.Mutex
}

func (db *Sqlite3Db) Open(fp string) {
	var err error

	db.Db, err = sqlite.OpenConn(fp, 0)
	if err != nil {
		DbL.Fatalln(err)
	}

	err = sqlitex.ExecScript(db.Db, `create table if not exists torrent (infohash text primary key,started boolean,addedat text,startedat text);`)

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

	serr := sqlitex.Exec(
		db.Db, `select 1 from torrent where infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			ret = stmt.ColumnInt(0) == 1
			return nil
		}, ih.HexString())
	if serr != nil {
		return false
	}
	return
}

func (db *Sqlite3Db) IsLocked(ih string) (ret bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(
		db.Db, `select locked from torrent where infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			ret = stmt.ColumnInt(0) != 0
			return nil
		}, ih)

	return
}

func (db *Sqlite3Db) HasStarted(ih string) (ret bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(
		db.Db, `select started from torrent where infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			ret = stmt.ColumnInt(0) != 0
			return nil
		}, ih)

	return
}

func (db *Sqlite3Db) SetLocked(ih string, b bool) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `update torrent set locked=? where infohash=?;`, nil, b, ih)
	return
}

func (db *Sqlite3Db) Add(ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	tn := time.Now().Format(time.RFC3339)
	err = sqlitex.Exec(db.Db, `insert into torrent (infohash,started,addedat,startedat) values (?,?,?,?) on conflict (infohash) do update set startedat=?;`, nil, ih.HexString(), 0, tn, tn, tn)
	return
}

func (db *Sqlite3Db) Delete(ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `delete from torrent where infohash=?;`, nil, ih.HexString())
	return
}

func (db *Sqlite3Db) Start(ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	tn := time.Now().Format(time.RFC3339)
	err = sqlitex.Exec(db.Db, `update torrent set started=?,startedat=? where infohash=?;`, nil, 1, tn, ih.HexString())
	return
}

func (db *Sqlite3Db) SetStarted(ih metainfo.Hash, inp bool) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `update torrent set started=? where infohash=?;`, nil, inp, ih.HexString())
	return
}

func (db *Sqlite3Db) GetTorrent(ih metainfo.Hash) (*Torrent, error) {
	var trnt Torrent
	var exists bool
	var serr error
	var terr error

	db.mu.Lock()
	defer db.mu.Unlock()
	serr = sqlitex.Exec(
		db.Db, `select * from torrent where infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			exists = true
			trnt.Infohash = ih
			trnt.Started = stmt.ColumnInt(1) != 0
			trnt.AddedAt, terr = time.Parse(time.RFC3339, stmt.GetText("addedat"))
			if terr != nil {
				DbL.Println(terr)
				return terr
			}
			trnt.StartedAt, terr = time.Parse(time.RFC3339, stmt.GetText("startedat"))
			if terr != nil {
				DbL.Println(terr)
				return terr
			}
			return nil
		}, ih.HexString())

	if serr != nil {
		return nil, serr
	}
	if !exists {
		return nil, fmt.Errorf("Torrent doesn't exist")
	}
	return &trnt, nil
}

func (db *Sqlite3Db) GetTorrents() (Trnts []*Torrent, err error) {
	Trnts = make([]*Torrent, 0)

	var serr error
	var terr error
	db.mu.Lock()
	defer db.mu.Unlock()
	serr = sqlitex.Exec(
		db.Db, `select * from torrent;`,
		func(stmt *sqlite.Stmt) error {
			var trnt Torrent
			trnt.Infohash, terr = MetafromHex(stmt.GetText("infohash"))
			if terr != nil {
				DbL.Println(terr)
				return terr
			}
			trnt.Started = stmt.ColumnInt(1) != 0
			trnt.AddedAt, terr = time.Parse(time.RFC3339, stmt.GetText("addedat"))
			if terr != nil {
				DbL.Println(terr)
				return terr
			}
			trnt.StartedAt, terr = time.Parse(time.RFC3339, stmt.GetText("startedat"))
			if terr != nil {
				DbL.Println(terr)
				return terr
			}
			Trnts = append(Trnts, &trnt)
			return nil
		})
	if serr != nil {
		return Trnts, serr
	}

	return Trnts, nil
}
