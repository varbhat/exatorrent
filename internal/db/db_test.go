package db

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/uptrace/bun"
	"log"
)

// init testdata here

var db *bun.DB
var ih1 metainfo.Hash
var ih2 metainfo.Hash

func init() {
	var err error

	//dsn := "postgres://postgres:postgres@localhost:5432/exatorrent?sslmode=disable"
	//db, err = InitDb("postgres", dsn)
	db, err = InitDb(Sqlite, "data.db")
	if err != nil {
		log.Fatal(err)
	}

	ih1, err = MetafromHex("2b9d4e7ab1b301c2b06eb469f62b6651d0757b94")
	if err != nil {
		log.Fatal(err)
	}

	ih2, err = MetafromHex("d4b9f6b77df092cef483a343b036a00e03416756")
	if err != nil {
		log.Fatal(err)
	}
}
