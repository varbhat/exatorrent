package db

import (
	"sync"

	"github.com/anacrolix/torrent/metainfo"
	sqlite "github.com/go-llsqlite/crawshaw"
	"github.com/go-llsqlite/crawshaw/sqlitex"
)

type SqliteLSDb struct {
	Db *sqlite.Conn
	mu sync.Mutex
}

func (db *SqliteLSDb) Open(fp string) {
	var err error

	db.Db, err = sqlite.OpenConn(fp, 0)
	if err != nil {
		DbL.Fatalln(err)
	}

	err = sqlitex.ExecScript(db.Db, `create table if not exists lockstatedb (infohash text primary key);`)

	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteLSDb) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.Db.Close()
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteLSDb) Lock(m metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `insert into lockstatedb (infohash) values (?) on conflict (infohash) do nothing;`, nil, m.HexString())
	return
}

func (db *SqliteLSDb) Unlock(m metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `delete from lockstatedb where infohash=?;`, nil, m.HexString())
	return
}

func (db *SqliteLSDb) IsLocked(m string) (b bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if sqlitex.Exec(
		db.Db, `select 1 from lockstatedb where infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			b = true
			return nil
		}, m) != nil {
		return false
	}
	return
}
