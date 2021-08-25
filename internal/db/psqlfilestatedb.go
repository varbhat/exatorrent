package db

import (
	"context"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PsqlFsDb struct {
	Db *pgxpool.Pool
}

func (db *PsqlFsDb) Open(dburl string) {
	var err error
	db.Db, err = pgxpool.Connect(context.Background(), dburl)
	if err != nil {
		DbL.Fatalln(err)
	}

	_, err = db.Db.Exec(context.Background(), `create table if not exists filestatedb (filepath text,infohash text, unique(filepath, infohash));`)
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *PsqlFsDb) Close() {
	db.Db.Close()
}

func (db *PsqlFsDb) Add(fp string, ih metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `insert into filestatedb (filepath,infohash) values ($1,$2) on conflict (filepath,infohash) do nothing;`, fp, ih.HexString())
	return
}

func (db *PsqlFsDb) Get(ih metainfo.Hash) (ret []string) {
	ret = make([]string, 0)
	rows, err := db.Db.Query(context.Background(), `select filepath from filestatedb WHERE infohash=$1;`, ih.HexString())
	if err != nil {
		DbL.Println(err)
		return
	}

	for rows.Next() {
		var fpstring string
		err = rows.Scan(&fpstring)
		if err != nil {
			DbL.Println(err)
			return
		}
		ret = append(ret, fpstring)
	}

	return
}

func (db *PsqlFsDb) Deletefile(fp string, ih metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `delete from filestatedb where filepath=$1 and infohash=$2;`, fp, ih.HexString())
	return err
}

func (db *PsqlFsDb) Delete(ih metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `delete from filestatedb where infohash=$1;`, ih.HexString())
	return
}
