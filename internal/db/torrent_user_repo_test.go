package db

import (
	"github.com/anacrolix/torrent/metainfo"
	"testing"
)

func TestTorrentUserRepo_Open(t *testing.T) {
	tur := NewTorrentUserRepo(db)

	tur.Open("")
}

func TestTorrentUserRepo_Add(t *testing.T) {
	tur := NewTorrentUserRepo(db)

	var err error
	for _, username := range []string{"evrins", "admin"} {
		for _, hash := range []metainfo.Hash{ih1, ih2} {
			err = tur.Add(username, hash)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestTorrentUserRepo_HasUser(t *testing.T) {
	tur := NewTorrentUserRepo(db)

	username := "admin"
	hasUser := tur.HasUser(username, ih1.HexString())
	t.Logf("hasUser: %v", hasUser)
}

func TestTorrentUserRepo_ListTorrents(t *testing.T) {
	tur := NewTorrentUserRepo(db)

	username := "admin"
	torrentList := tur.ListTorrents(username)
	for _, hash := range torrentList {
		t.Logf("%v", hash)
	}
}

func TestTorrentUserRepo_ListUsers(t *testing.T) {
	tur := NewTorrentUserRepo(db)

	userList := tur.ListUsers(ih1)
	t.Logf("%+v", userList)
}

func TestTorrentUserRepo_Remove(t *testing.T) {
	tur := NewTorrentUserRepo(db)

	err := tur.Remove("admin", ih1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTorrentUserRepo_RemoveAll(t *testing.T) {
	tur := NewTorrentUserRepo(db)

	err := tur.RemoveAll("evrins")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTorrentUserRepo_RemoveAllMi(t *testing.T) {
	tur := NewTorrentUserRepo(db)

	err := tur.RemoveAllMi(ih1)
	if err != nil {
		t.Fatal(err)
	}
}
