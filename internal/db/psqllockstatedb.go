package db

import (
	"context"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PsqlLsDb struct {
	Db *pgxpool.Pool
}

func (db *PsqlLsDb) Open(dburl string) {
	var err error
	db.Db, err = pgxpool.Connect(context.Background(), dburl)
	if err != nil {
		DbL.Fatalln(err)
	}

	_, err = db.Db.Exec(context.Background(), `create table if not exists lockstatedb (infohash text primary key);`)
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *PsqlLsDb) Close() {
	db.Db.Close()
}

func (db *PsqlLsDb) Lock(m metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `insert into lockstatedb (infohash) values ($1) on conflict (infohash) do nothing;`, m.HexString())
	return
}

func (db *PsqlLsDb) Unlock(m metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `delete from lockstatedb where infohash=$1;`, m.HexString())
	return
}

func (db *PsqlLsDb) IsLocked(m string) (b bool) {
	if db.Db.QueryRow(context.Background(), `select true from lockstatedb where infohash=$1;`, m).Scan(&b) != nil {
		return false
	}
	return
}
