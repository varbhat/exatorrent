package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PsqlTDb struct {
	Db *pgxpool.Pool
}

func (db *PsqlTDb) Open(dburl string) {
	var err error
	db.Db, err = pgxpool.Connect(context.Background(), dburl)
	if err != nil {
		DbL.Fatalln(err)
	}

	_, err = db.Db.Exec(context.Background(), `create table if not exists trackerdb (url text primary key);`)
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *PsqlTDb) Close() {
	db.Db.Close()
}

func (db *PsqlTDb) Add(url string) {
	_, _ = db.Db.Exec(context.Background(), `insert into trackerdb (url) values ($1) on conflict (url) do nothing;`, url)
}

func (db *PsqlTDb) Delete(url string) {
	_, _ = db.Db.Exec(context.Background(), `delete from trackerdb where url=$1;`, url)
}

func (db *PsqlTDb) DeleteN(count int) {
	_, _ = db.Db.Exec(context.Background(), `delete from trackerdb where url in (select url from trackerdb limit $1);`, count)
}

func (db *PsqlTDb) DeleteAll() {
	_, _ = db.Db.Exec(context.Background(), `delete from trackerdb;`)
}

func (db *PsqlTDb) Count() (ret int64) {
	row := db.Db.QueryRow(context.Background(), `select count(*) from trackerdb;`)
	_ = row.Scan(&ret)
	return
}

func (db *PsqlTDb) Get() (ret []string) {
	rows, err := db.Db.Query(context.Background(), `select url from trackerdb;`)
	if err != nil {
		DbL.Println(err)
	}

	for rows.Next() {
		var tr string
		err := rows.Scan(&tr)
		if err != nil {
			DbL.Println(err)
			return
		}
		ret = append(ret, tr)
	}
	return
}
