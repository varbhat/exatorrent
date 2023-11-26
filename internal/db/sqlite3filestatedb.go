package db

import (
	"sync"

	"github.com/anacrolix/torrent/metainfo"
	sqlite "github.com/go-llsqlite/crawshaw"
	"github.com/go-llsqlite/crawshaw/sqlitex"
)

type SqliteFSDb struct {
	Db *sqlite.Conn
	mu sync.Mutex
}

func (db *SqliteFSDb) Open(fp string) {
	var err error

	db.Db, err = sqlite.OpenConn(fp, 0)
	if err != nil {
		DbL.Fatalln(err)
	}

	err = sqlitex.ExecScript(db.Db, `create table if not exists filestatedb (filepath text,infohash text, unique(filepath, infohash));`)

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
	err = sqlitex.Exec(db.Db, `insert into filestatedb (filepath,infohash) values (?,?) on conflict (filepath,infohash) do nothing;`, nil, fp, ih.HexString())
	return
}

func (db *SqliteFSDb) Get(ih metainfo.Hash) (ret []string) {
	ret = make([]string, 0)
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(
		db.Db, `select filepath from filestatedb where infohash=?;`,
		func(stmt *sqlite.Stmt) error {
			ret = append(ret, stmt.GetText("filepath"))
			return nil
		}, ih.HexString())

	return
}

func (db *SqliteFSDb) Delete(ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `delete from filestatedb where infohash=?;`, nil, ih.HexString())
	return
}

func (db *SqliteFSDb) Deletefile(fp string, ih metainfo.Hash) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	err = sqlitex.Exec(db.Db, `delete from filestatedb where filepath=? and infohash=?;`, nil, fp, ih.HexString())
	return err
}
