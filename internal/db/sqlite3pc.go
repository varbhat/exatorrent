//go:build cgo && !nosqlite
// +build cgo,!nosqlite

package db

import (
	"path/filepath"
	"sync"

	sqlite "github.com/go-llsqlite/crawshaw"
	"github.com/go-llsqlite/crawshaw/sqlitex"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

type sqlitePieceCompletion struct {
	mu sync.Mutex
	db *sqlite.Conn
}

func NewSqlitePieceCompletion(dir string) (ret *sqlitePieceCompletion, err error) {
	p := filepath.Join(dir, "pcomp.db")
	db, err := sqlite.OpenConn(p, 0)
	if err != nil {
		return
	}
	err = sqlitex.ExecScript(db, `create table if not exists pcomp (infohash text, pindex integer, complete boolean, unique(infohash, pindex));`)
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
	err = sqlitex.Exec(
		me.db, `select complete from pcomp where infohash=? and pindex=?;`,
		func(stmt *sqlite.Stmt) error {
			c.Complete = stmt.ColumnInt(0) != 0
			c.Ok = true
			return nil
		},
		pk.InfoHash.HexString(), pk.Index)
	return
}

func (me *sqlitePieceCompletion) Set(pk metainfo.PieceKey, b bool) error {
	me.mu.Lock()
	defer me.mu.Unlock()
	return sqlitex.Exec(
		me.db,
		`insert into pcomp (infohash,pindex,complete) values (?,?,?) on conflict (infohash,pindex) do update set complete=?;`,
		nil,
		pk.InfoHash.HexString(), pk.Index, b, b)
}

func (me *sqlitePieceCompletion) Delete(m metainfo.Hash) {
	me.mu.Lock()
	defer me.mu.Unlock()
	if me.db != nil {
		_ = sqlitex.Exec(
			me.db,
			`delete from pcomp where infohash=?;`,
			nil,
			m.HexString())
	}
}

func (me *sqlitePieceCompletion) Close() (err error) {
	me.mu.Lock()
	defer me.mu.Unlock()
	if me.db != nil {
		err = me.db.Close()
		me.db = nil
	}
	return
}
