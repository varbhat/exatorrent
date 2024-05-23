//go:build cgo
// +build cgo

package db

import (
	"database/sql"
	"sync"
)

type SqliteTdb struct {
	Db *sql.DB
	mu sync.Mutex
}

func (db *SqliteTdb) Open(fp string) {
	var err error

	db.Db, err = sql.Open("sqlite3", fp)
	if err != nil {
		DbL.Fatalln(err)
	}
	db.Db.SetMaxOpenConns(1)

	_, err = db.Db.Exec(`create table if not exists trackerdb (url text primary key);`)

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
	_, err := db.Db.Exec(`insert into trackerdb (url) values (?);`, url)
	if err != nil {
		DbL.Println(err)
	}
}

func (db *SqliteTdb) Delete(url string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err := db.Db.Exec(`delete from trackerdb where url=?;`, url)
	if err != nil {
		DbL.Println(err)
	}
}

func (db *SqliteTdb) DeleteN(count int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err := db.Db.Exec(`delete from trackerdb where url in (select url from trackerdb limit ?);`, count)
	if err != nil {
		DbL.Println(err)
	}
}

func (db *SqliteTdb) DeleteAll() {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err := db.Db.Exec(`delete from trackerdb;`)
	if err != nil {
		DbL.Println(err)
	}
}

func (db *SqliteTdb) Count() (ret int64) {
	db.mu.Lock()
	defer db.mu.Unlock()
	row := db.Db.QueryRow(`select count(*) from trackerdb;`)
	err := row.Err()
	if err != nil {
		DbL.Println(err)
		return
	}
	err = row.Scan(&ret)
	if err != nil {
		return
	}
	return
}

func (db *SqliteTdb) Get() (ret []string) {
	ret = make([]string, 0)
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.Db.Query("select url from trackerdb;")
	if err != nil {
		return
	}
	var url string
	for rows.Next() {
		err = rows.Scan(&url)
		ret = append(ret, url)
	}
	return
}
