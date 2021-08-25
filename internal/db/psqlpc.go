package db

import (
	"context"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

type psqlPieceCompletion struct {
	db *pgxpool.Pool
}

func NewPsqlPieceCompletion(databaseurl string) (ret *psqlPieceCompletion, err error) {
	db, err := pgxpool.Connect(context.Background(), databaseurl)
	if err != nil {
		return
	}

	_, err = db.Exec(context.Background(), `create table if not exists pcomp (infohash text,pindex integer,complete boolean, unique(infohash, pindex));`)
	if err != nil {
		db.Close()
		return
	}

	ret = &psqlPieceCompletion{db: db}
	return
}

func (me *psqlPieceCompletion) Get(pk metainfo.PieceKey) (c storage.Completion, err error) {
	if me.db.QueryRow(context.Background(), `select complete from pcomp where infohash=$1 and pindex=$2;`, pk.InfoHash.HexString(), pk.Index).Scan(&c.Complete) == nil {
		c.Ok = true
	}
	return
}

func (me *psqlPieceCompletion) Set(pk metainfo.PieceKey, b bool) (err error) {
	_, err = me.db.Exec(context.Background(), `insert into pcomp (infohash,pindex,complete) values ($1,$2,$3) on conflict (infohash,pindex) do update set complete=$3;`, pk.InfoHash.HexString(), pk.Index, b)
	return
}

func (me *psqlPieceCompletion) Delete(m metainfo.Hash) {
	if me != nil {
		_, _ = me.db.Exec(context.Background(), `delete from pcomp where infohash=$1;`, m.HexString())
	}
}

func (me *psqlPieceCompletion) Close() error {
	me.db.Close()
	return nil
}
