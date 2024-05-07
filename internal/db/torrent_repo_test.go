package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTorrentRepo_Open(t *testing.T) {
	tr := NewTorrentRepo(db)

	tr.Open("")
}

func TestTorrentRepo_Add(t *testing.T) {
	tr := NewTorrentRepo(db)

	err := tr.Add(ih1)
	if err != nil {
		t.Fatal(err)
	}

	err = tr.Add(ih2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTorrentRepo_Delete(t *testing.T) {
	tr := NewTorrentRepo(db)

	err := tr.Delete(ih1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTorrentRepo_Exists(t *testing.T) {
	tr := NewTorrentRepo(db)

	exists := tr.Exists(ih1)
	t.Log(exists)
}

func TestTorrentRepo_GetTorrent(t *testing.T) {
	tr := NewTorrentRepo(db)

	torrent, err := tr.GetTorrent(ih1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(torrent)
}

func TestTorrentRepo_GetTorrents(t *testing.T) {
	tr := NewTorrentRepo(db)

	torrentList, err := tr.GetTorrents()
	if err != nil {
		t.Fatal(err)
	}
	for _, torrent := range torrentList {
		t.Logf("%+v", torrent)
	}
}

func TestTorrentRepo_HasStarted(t *testing.T) {
	tr := NewTorrentRepo(db)

	hash := "12"
	hasStarted := tr.HasStarted(hash)
	t.Logf("hasStarted: %v", hasStarted)
}

func TestTorrentRepo_SetStarted(t *testing.T) {
	tr := NewTorrentRepo(db)

	err := tr.SetStarted(ih1, false)
	if err != nil {
		t.Fatal(err)
	}

	started := tr.HasStarted(ih1.HexString())
	assert.Equal(t, false, started, "")

	err = tr.SetStarted(ih1, true)
	if err != nil {
		t.Fatal(err)
	}

	started = tr.HasStarted(ih1.HexString())
	assert.Equal(t, true, started, "")
}

func TestTorrentRepo_Start(t *testing.T) {
	tr := NewTorrentRepo(db)

	err := tr.Start(ih1)
	if err != nil {
		t.Fatal(err)
	}

	torrent, err := tr.GetTorrent(ih1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(torrent)
}
