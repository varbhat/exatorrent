package db

import (
	"database/sql"
	"sync"

	"github.com/anacrolix/torrent/metainfo"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteFSDb struct {
	Db *sql.DB
	mu sync.Mutex
}

func (db *SqliteFSDb) Open(fp string) {
	var err error

	db.Db, err = sql.Open("sqlite3", fp)
	if err != nil {
		DbL.Fatalln(err)
	}
	db.Db.SetMaxOpenConns(1)
	_, err = db.Db.Exec(`create table if not exists filestatedb (filepath text,infohash text, unique(filepath, infohash));`)

	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteFSDb) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.Db.Close()
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteFSDb) Add(fp string, ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`insert into filestatedb (filepath,infohash) values (?,?) on conflict (filepath,infohash) do nothing;`, fp, ih.HexString())
	return
}

func (db *SqliteFSDb) Get(ih metainfo.Hash) (ret []string) {
	ret = make([]string, 0)
	db.mu.Lock()
	defer db.mu.Unlock()
	rows, err := db.Db.Query(`select filepath from filestatedb where infohash=?;`, ih.HexString())

	if err != nil {
		DbL.Printf("fail to query filestate. err: %v", err)
		return
	}

	var fpstring string
	for rows.Next() {
		err = rows.Scan(&fpstring)
		if err != nil {
			DbL.Println(err)
			return
		}
		ret = append(ret, fpstring)
	}
	return
}

func (db *SqliteFSDb) Delete(ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`delete from filestatedb where infohash=?;`, ih.HexString())
	return
}

func (db *SqliteFSDb) Deletefile(fp string, ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`delete from filestatedb where filepath=? and infohash=?;`, fp, ih.HexString())
	return err
}
