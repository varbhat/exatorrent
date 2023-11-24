//go:build cgo
// +build cgo

package db

import (
	"sync"

	sqlite "github.com/go-llsqlite/crawshaw"
	"github.com/go-llsqlite/crawshaw/sqlitex"
)

type SqliteTdb struct {
	Db *sqlite.Conn
	mu sync.Mutex
}

func (db *SqliteTdb) Open(fp string) {
	var err error

	db.Db, err = sqlite.OpenConn(fp, 0)
	if err != nil {
		DbL.Fatalln(err)
	}

	err = sqlitex.ExecScript(db.Db, `create table if not exists trackerdb (url text primary key);`)

	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteTdb) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.Db.Close()
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *SqliteTdb) Add(url string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(db.Db, `insert into trackerdb (url) values (?);`, nil, url)
}

func (db *SqliteTdb) Delete(url string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(db.Db, `delete from trackerdb where url=?;`, nil, url)
}

func (db *SqliteTdb) DeleteN(count int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(db.Db, `delete from trackerdb where url in (select url from trackerdb limit ?);`, nil, count)
}

func (db *SqliteTdb) DeleteAll() {
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(db.Db, `delete from trackerdb;`, nil)
}

func (db *SqliteTdb) Count() (ret int64) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(
		db.Db, `select count(*) from trackerdb;`,
		func(stmt *sqlite.Stmt) error {
			ret = stmt.ColumnInt64(0)
			return nil
		})
	return
}

func (db *SqliteTdb) Get() (ret []string) {
	ret = make([]string, 0)
	db.mu.Lock()
	defer db.mu.Unlock()
	_ = sqlitex.Exec(
		db.Db, `select url from trackerdb;`,
		func(stmt *sqlite.Stmt) error {
			ret = append(ret, stmt.GetText("url"))
			return nil
		})

	return
}
