//go:build !cgo
// +build !cgo

package core

import (
	"github.com/anacrolix/torrent"
)

func sqliteSetup(tc *torrent.ClientConfig) {
	Err.Fatalln("Postgresql Connection URL was not provided")
}
