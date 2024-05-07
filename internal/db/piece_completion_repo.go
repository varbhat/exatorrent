package db

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/uptrace/bun"
)

type PieceCompletion struct {
	bun.BaseModel `bun:"table:pcomp,alias:pc"`

	InfoHash string `bun:"infohash,unique:pc_uniq_key"`
	PIndex   int    `bun:"pindex,unique:pc_uniq_key"`
	Complete bool   `bun:"complete"`
}

type PieceCompletionRepo struct {
	conn *bun.DB
}

func NewPieceCompletionRepo(conn *bun.DB) *PieceCompletionRepo {
	return &PieceCompletionRepo{conn: conn}
}

func (pcr *PieceCompletionRepo) Open() {
	_, err := pcr.conn.
		NewCreateTable().
		Model(&PieceCompletion{}).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		DbL.Fatalf("failed creating table for PieceCompletion: %v", err)
	}
}

func (pcr *PieceCompletionRepo) Get(pk metainfo.PieceKey) (c storage.Completion, err error) {
	err = pcr.conn.
		NewSelect().
		Model(&PieceCompletion{}).
		Column("complete").
		Where("infohash = ?", pk.InfoHash.HexString()).
		Where("pindex = ?", pk.Index).
		Scan(context.Background(), &c.Complete)

	if err == nil {
		c.Ok = true
	}
	return c, err
}

func (pcr *PieceCompletionRepo) Set(pk metainfo.PieceKey, complete bool) error {
	_, err := pcr.conn.
		NewInsert().
		Model(&PieceCompletion{
			InfoHash: pk.InfoHash.HexString(),
			PIndex:   pk.Index,
			Complete: complete,
		}).
		On("conflict (infohash, pindex) do update").
		Set("complete = EXCLUDED.complete").
		Exec(context.Background())
	return err
}

func (pcr *PieceCompletionRepo) Close() error {
	return pcr.conn.Close()
}

func (pcr *PieceCompletionRepo) Delete(hash metainfo.Hash) {
	_, err := pcr.conn.
		NewDelete().
		Model(&PieceCompletion{}).
		Where("infohash = ?", hash.HexString()).
		Exec(context.Background())

	if err != nil {
		DbL.Fatalf("failed deleting PieceCompletion: %v", err)
	}
}
