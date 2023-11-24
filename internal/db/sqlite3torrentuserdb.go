//go:build cgo
// +build cgo

package db

import (
	"sync"

	sqlite "github.com/go-llsqlite/crawshaw"
	"github.com/go-llsqlite/crawshaw/sqlitex"

	"github.com/anacrolix/torrent/metainfo"
)

type SqliteTorrentUserDb struct {
	Db *sqlite.Conn
	mu sync.Mutex
}

func (db *SqliteTorrentUserDb) Open(fp string) {
	var err error

	db.Db, err = sqlite.OpenConn(fp, 0)
	if err != nil {
		DbL.Fatalln(err)
	}

	err = sqlitex.ExecScript(db.Db, `create table if not exists torrentuserdb (username text,infohash text, unique(username,infohash));`)

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
	err = sqlitex.Exec(db.Db, `insert into torrentuserdb (username,infohash) values (?,?) on conflict (username,infohash) do nothing;`, nil, username, ih.HexString())
	return
}

func (db *SqliteTorrentUserDb) Remove(username string, ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `delete from torrentuserdb where username=? and infohash=?;`, nil, username, ih.HexString())
	return
}

func (db *SqliteTorrentUserDb) RemoveAll(username string) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `delete from torrentuserdb where username=?;`, nil, username)
	return
}

func (db *SqliteTorrentUserDb) RemoveAllMi(mi metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `delete from torrentuserdb where infohash=?;`, nil, mi.HexString())
	return
}

func (db *SqliteTorrentUserDb) HasUser(username string, ih string) (ret bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	serr := sqlitex.Exec(
		db.Db, `select 1 from torrentuserdb where username=? and infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			ret = stmt.ColumnInt(0) == 1
			return nil
		}, username, ih)
	if serr != nil {
		return false
	}
	return
}

func (db *SqliteTorrentUserDb) ListTorrents(username string) (ret []metainfo.Hash) {
	ret = make([]metainfo.Hash, 0)
	db.mu.Lock()
	defer db.mu.Unlock()

	_ = sqlitex.Exec(
		db.Db, `select infohash from torrentuserdb where username=?;`,
		func(stmt *sqlite.Stmt) error {
			tm, terr := MetafromHex(stmt.GetText("infohash"))
			if terr != nil {
				DbL.Println(terr)
				return terr
			}
			ret = append(ret, tm)
			return nil
		}, username)
	return
}

func (db *SqliteTorrentUserDb) ListUsers(mi metainfo.Hash) (ret []string) {
	ret = make([]string, 0)
	db.mu.Lock()
	defer db.mu.Unlock()

	_ = sqlitex.Exec(
		db.Db, `select username from torrentuserdb where infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			username := stmt.GetText("username")
			ret = append(ret, username)
			return nil
		}, mi.HexString())
	return
}
