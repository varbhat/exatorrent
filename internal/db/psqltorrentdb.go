package db

import (
	"context"
	"time"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PsqlTrntDb struct {
	Db *pgxpool.Pool
}

func (db *PsqlTrntDb) Open(dburl string) {
	var err error
	db.Db, err = pgxpool.Connect(context.Background(), dburl)
	if err != nil {
		DbL.Fatalln(err)
	}
	_, err = db.Db.Exec(context.Background(), `create table if not exists torrent (infohash text primary key,started boolean,addedat timestamptz,startedat timestamptz);`)
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *PsqlTrntDb) Close() {
	db.Db.Close()
}

func (db *PsqlTrntDb) Exists(ih metainfo.Hash) (ret bool) {
	var err error
	row := db.Db.QueryRow(context.Background(), `select true from torrent where infohash=$1;`, ih.HexString())
	err = row.Scan(&ret)
	if err != nil {
		return false
	}
	return
}

func (db *PsqlTrntDb) HasStarted(ih string) (ret bool) {
	row := db.Db.QueryRow(context.Background(), `select started from torrent where infohash=$1;`, ih)
	err := row.Scan(&ret)
	if err != nil {
		return false
	}
	return
}

func (db *PsqlTrntDb) Add(ih metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `insert into torrent (infohash,started,addedat,startedat) values ($1,$2,$3,$4) on conflict (infohash) do update set startedat=$4;`, ih.HexString(), false, time.Now(), time.Now())
	return
}

func (db *PsqlTrntDb) Delete(ih metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `delete from torrent where infohash=$1;`, ih.HexString())
	return
}

func (db *PsqlTrntDb) Start(ih metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `update torrent set started=$1,startedat=$2 where infohash=$3;`, true, time.Now(), ih.HexString())
	return
}

func (db *PsqlTrntDb) SetStarted(ih metainfo.Hash, inp bool) (err error) {
	_, err = db.Db.Exec(context.Background(), `update torrent set started=$1 where infohash=$2;`, inp, ih.HexString())
	return
}

func (db *PsqlTrntDb) GetTorrent(ih metainfo.Hash) (*Torrent, error) {
	var trnt Torrent
	var infoh string
	row := db.Db.QueryRow(context.Background(), `select * from torrent where infohash=$1;`, ih.HexString())
	err := row.Scan(&infoh, &trnt.Started, &trnt.AddedAt, &trnt.StartedAt)
	if err != nil {
		return nil, err
	}
	trnt.Infohash = ih
	return &trnt, nil
}

func (db *PsqlTrntDb) GetTorrents() (Trnts []*Torrent, err error) {
	Trnts = make([]*Torrent, 0)
	rows, err := db.Db.Query(context.Background(), `select * from torrent;`)
	if err != nil {
		DbL.Println(err)
		return
	}

	for rows.Next() {
		var trnt Torrent
		var ih string
		err = rows.Scan(&ih, &trnt.Started, &trnt.AddedAt, &trnt.StartedAt)
		if err != nil {
			DbL.Println(err)
			return Trnts, err
		}

		trnt.Infohash, err = MetafromHex(ih)
		if err != nil {
			DbL.Println(err)

			return Trnts, err
		}
		Trnts = append(Trnts, &trnt)
	}

	return Trnts, rows.Err()
}
