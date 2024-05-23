package db

import (
	"database/sql"
	"sync"

	"github.com/anacrolix/torrent/metainfo"
)

type SqliteLSDb struct {
	Db *sql.DB
	mu sync.Mutex
}

func (db *SqliteLSDb) Open(fp string) {
	var err error

	db.Db, err = sql.Open("sqlite3", fp)
	db.Db.SetMaxOpenConns(1)
	if err != nil {
		DbL.Fatalln(err)
	}

	_, err = db.Db.Exec(`create table if not exists lockstatedb (infohash text primary key);`)

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
	_, err = db.Db.Exec(`insert into lockstatedb (infohash) values (?) on conflict (infohash) do nothing;`, m.HexString())
	return
}

func (db *SqliteLSDb) Unlock(m metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`delete from lockstatedb where infohash=?;`, m.HexString())
	return
}

func (db *SqliteLSDb) IsLocked(m string) (b bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.Db.QueryRow(`select exists ( select 1 from lockstatedb where infohash=?);`, m)
	err := row.Err()
	if err != nil {
		return false
	}

	err = row.Scan(&b)
	if err != nil {
		DbL.Printf("failed to check lockstatedb infohash %s: %s", m, err)
		return false
	}

	return b

}
