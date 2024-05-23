//go:build cgo
// +build cgo

package db

import (
	"database/sql"
	"sync"

	"github.com/anacrolix/torrent/metainfo"
)

type SqliteTorrentUserDb struct {
	Db *sql.DB
	mu sync.Mutex
}

func (db *SqliteTorrentUserDb) Open(fp string) {
	var err error

	db.Db, err = sql.Open("sqlite3", fp)
	if err != nil {
		DbL.Fatalln(err)
	}
	db.Db.SetMaxOpenConns(1)
	_, err = db.Db.Exec(`create table if not exists torrentuserdb (username text,infohash text, unique(username,infohash));`)

	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteTorrentUserDb) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.Db.Close()
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteTorrentUserDb) Add(username string, ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`insert into torrentuserdb (username,infohash) values (?,?) on conflict (username,infohash) do nothing;`, username, ih.HexString())
	return
}

func (db *SqliteTorrentUserDb) Remove(username string, ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`delete from torrentuserdb where username=? and infohash=?;`, username, ih.HexString())
	return
}

func (db *SqliteTorrentUserDb) RemoveAll(username string) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`delete from torrentuserdb where username=?;`, username)
	return
}

func (db *SqliteTorrentUserDb) RemoveAllMi(mi metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`delete from torrentuserdb where infohash=?;`, mi.HexString())
	return
}

func (db *SqliteTorrentUserDb) HasUser(username string, ih string) (ret bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.Db.QueryRow(`select exists (select 1 from torrentuserdb where username=? and infohash=?);`, ih)
	var err error
	err = row.Err()
	if err != nil {
		return
	}
	err = row.Scan(&ret)
	if err != nil {
		return
	}

	return
}

func (db *SqliteTorrentUserDb) ListTorrents(username string) (ret []metainfo.Hash) {
	ret = make([]metainfo.Hash, 0)
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.Db.Query(`select infohash from torrentuserdb where username=?;`, username)
	if err != nil {
		return
	}
	var hash string
	var ih metainfo.Hash
	for rows.Next() {
		err = rows.Scan(&hash)
		if err != nil {
			return
		}
		ih, err = MetafromHex(hash)
		if err != nil {
			return
		}
		ret = append(ret, ih)
	}
	return
}

func (db *SqliteTorrentUserDb) ListUsers(mi metainfo.Hash) (ret []string) {
	ret = make([]string, 0)
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.Db.Query(`select username from torrentuserdb where infohash=?;`, mi.HexString())
	if err != nil {
		return
	}
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return
		}
		ret = append(ret, username)
	}

	return
}
