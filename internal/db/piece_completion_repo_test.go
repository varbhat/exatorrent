package db

import (
	"github.com/anacrolix/torrent/metainfo"
	"testing"
)

func TestPieceCompletionRepo_Open(t *testing.T) {
	pcr := NewPieceCompletionRepo(db)

	pcr.Open()
}

func TestPieceCompletionRepo_Delete(t *testing.T) {
	pcr := NewPieceCompletionRepo(db)

	pcr.Delete(ih1)
}

func TestPieceCompletionRepo_Get(t *testing.T) {
	pcr := NewPieceCompletionRepo(db)

	c, err := pcr.Get(metainfo.PieceKey{
		InfoHash: ih1,
		Index:    8,
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(c)
}

func TestPieceCompletionRepo_Set(t *testing.T) {
	pcr := NewPieceCompletionRepo(db)

	var err error

	for _, hash := range []metainfo.Hash{ih1, ih2} {
		for i := 0; i < 10; i++ {
			pk := metainfo.PieceKey{
				InfoHash: hash,
				Index:    i,
			}
			err = pcr.Set(pk, i%2 == 0)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
