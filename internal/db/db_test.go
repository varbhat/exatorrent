package db

import (
	"github.com/anacrolix/torrent/metainfo"
	"log"
)

var ih1 metainfo.Hash

func init() {
	var err error
	ih1, err = MetafromHex("c28e93bfbe5036b6d3612859e21836859b02c97b")
	if err != nil {
		log.Fatalf("failed to load metainfo: %v", err)
	}
}
