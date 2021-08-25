package db

import (
	"context"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PsqlTrntUserDb struct {
	Db *pgxpool.Pool
}

func (db *PsqlTrntUserDb) Open(dburl string) {
	var err error
	db.Db, err = pgxpool.Connect(context.Background(), dburl)
	if err != nil {
		DbL.Fatalln(err)
	}
	_, err = db.Db.Exec(context.Background(), `create table if not exists torrentuserdb (username text,infohash text, unique(username,infohash));`)
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *PsqlTrntUserDb) Close() {
	db.Db.Close()
}

func (db *PsqlTrntUserDb) Add(username string, ih metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `insert into torrentuserdb (username,infohash) values ($1,$2) on conflict (username,infohash) do nothing;`, username, ih.HexString())
	return
}

func (db *PsqlTrntUserDb) Remove(username string, ih metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `delete from torrentuserdb where username=$1 and infohash=$2;`, username, ih.HexString())
	return
}

func (db *PsqlTrntUserDb) RemoveAll(username string) (err error) {
	_, err = db.Db.Exec(context.Background(), `delete from torrentuserdb where username=$1;`, username)
	return
}

func (db *PsqlTrntUserDb) RemoveAllMi(mi metainfo.Hash) (err error) {
	_, err = db.Db.Exec(context.Background(), `delete from torrentuserdb where infohash=$1;`, mi.HexString())
	return
}

func (db *PsqlTrntUserDb) HasUser(username string, ih string) (ret bool) {
	var err error
	row := db.Db.QueryRow(context.Background(), `select true from torrentuserdb where username=$1 and infohash=$2;`, username, ih)
	err = row.Scan(&ret)
	if err != nil {
		return false
	}
	return
}

func (db *PsqlTrntUserDb) ListTorrents(username string) (ret []metainfo.Hash) {
	rows, err := db.Db.Query(context.Background(), `select infohash from torrentuserdb where username=$1;`, username)
	if err != nil {
		DbL.Println(err)
		return
	}

	for rows.Next() {
		var ih string
		err = rows.Scan(&ih)
		if err != nil {
			DbL.Println(err)
			return
		}

		infoh, err := MetafromHex(ih)
		if err != nil {
			DbL.Println(err)
			return
		}
		ret = append(ret, infoh)
	}

	return

}

func (db *PsqlTrntUserDb) ListUsers(mi metainfo.Hash) (ret []string) {
	ret = make([]string, 0)

	rows, err := db.Db.Query(context.Background(), `select username from torrentuserdb where infohash=$1;`, mi.HexString())
	if err != nil {
		DbL.Println(err)
		return
	}

	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			DbL.Println(err)
			return
		}

		ret = append(ret, username)
	}

	return

}
