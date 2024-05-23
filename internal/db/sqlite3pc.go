//go:build cgo && !nosqlite
// +build cgo,!nosqlite

package db

import (
	"database/sql"
	"path/filepath"
	"sync"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

type sqlitePieceCompletion struct {
	mu sync.Mutex
	db *sql.DB
}

func NewSqlitePieceCompletion(dir string) (ret *sqlitePieceCompletion, err error) {
	p := filepath.Join(dir, "pcomp.db")
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`create table if not exists pcomp (infohash text, pindex integer, complete boolean, unique(infohash, pindex));`)
	if err != nil {
		_ = db.Close()
		return
	}
	ret = &sqlitePieceCompletion{db: db}
	return
}

func (me *sqlitePieceCompletion) Get(pk metainfo.PieceKey) (c storage.Completion, err error) {
	me.mu.Lock()
	defer me.mu.Unlock()
	row := me.db.QueryRow(`select complete from pcomp where infohash=? and pindex=?;`, pk.InfoHash.HexString(), pk.Index)
	err = row.Err()
	if err != nil {
		c.Ok = false
		c.Err = err
		return
	}

	err = row.Scan(&c.Complete)
	if err != nil {
		c.Ok = false
		c.Err = err
		return
	}
	c.Ok = true
	return
}

func (me *sqlitePieceCompletion) Set(pk metainfo.PieceKey, b bool) error {
	me.mu.Lock()
	defer me.mu.Unlock()
	_, err := me.db.Exec(`insert into pcomp (infohash,pindex,complete) values (?,?,?) on conflict (infohash,pindex) do update set complete= excluded.complete;`,
		pk.InfoHash.HexString(), pk.Index, b)
	return err
}

func (me *sqlitePieceCompletion) Delete(m metainfo.Hash) {
	me.mu.Lock()
	defer me.mu.Unlock()
	if me.db == nil {
		return
	}
	_, err := me.db.Exec(`delete from pcomp where infohash=?;`, m.HexString())
	if err != nil {
		DbL.Printf("failed to delete pcomp infohash %s: %v", m.HexString(), err)
	}
}

func (me *sqlitePieceCompletion) Close() (err error) {
	me.mu.Lock()
	defer me.mu.Unlock()
	if me.db == nil {
		return nil
	}

	err = me.db.Close()
	me.db = nil
	return
}
