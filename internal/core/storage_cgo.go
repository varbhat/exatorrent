//go:build cgo
// +build cgo

package core

import (
	"path/filepath"

	"github.com/varbhat/exatorrent/internal/db"

	"github.com/anacrolix/torrent"
)

func sqliteSetup(tc *torrent.ClientConfig) {
	var err error

	Engine.TorDb = &db.Sqlite3Db{}
	Engine.TorDb.Open(filepath.Join(Dirconfig.DataDir, "torc.db"))

	Engine.TrackerDB = &db.SqliteTdb{}
	Engine.TrackerDB.Open(filepath.Join(Dirconfig.DataDir, "trackers.db"))

	Engine.FsDb = &db.SqliteFSDb{}
	Engine.FsDb.Open(filepath.Join(Dirconfig.DataDir, "filestate.db"))

	Engine.LsDb = &db.SqliteLSDb{}
	Engine.LsDb.Open(filepath.Join(Dirconfig.DataDir, "lockstate.db"))

	Engine.UDb = &db.Sqlite3UserDb{}
	Engine.UDb.Open(filepath.Join(Dirconfig.DataDir, "user.db"))

	Engine.TUDb = &db.SqliteTorrentUserDb{}
	Engine.TUDb.Open(filepath.Join(Dirconfig.DataDir, "torrentuser.db"))

	Engine.PcDb, err = db.NewSqlitePieceCompletion(Dirconfig.DataDir)

	if err != nil {
		Err.Fatalln("Unable to create sqlite3 database for PieceCompletion")
	}

}
