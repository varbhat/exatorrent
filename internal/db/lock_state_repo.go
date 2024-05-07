package db

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/uptrace/bun"
)

type LockState struct {
	bun.BaseModel `bun:"table:lockstatedb"`

	InfoHash string `bun:"infohash,pk"`
}

type LockStateRepo struct {
	conn *bun.DB
}

func NewLockStateRepo(db *bun.DB) *LockStateRepo {
	return &LockStateRepo{conn: db}
}

func (lsr *LockStateRepo) Open(s string) {
	var _, err = lsr.conn.NewCreateTable().
		Model(&LockState{}).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		DbL.Fatalf("failed to create table: %v", err)
	}
}

func (lsr *LockStateRepo) Close() {
	// do nothing
}

func (lsr *LockStateRepo) Lock(hash metainfo.Hash) error {
	_, err := lsr.conn.
		NewInsert().
		Model(&LockState{InfoHash: hash.HexString()}).
		On("conflict (infohash) do nothing").
		Exec(context.Background())
	return err
}

func (lsr *LockStateRepo) Unlock(hash metainfo.Hash) error {
	_, err := lsr.conn.
		NewDelete().
		Model(&LockState{}).
		Where("infohash = ?", hash.HexString()).
		Exec(context.Background())
	return err
}

func (lsr *LockStateRepo) IsLocked(s string) bool {
	result, err := lsr.conn.
		NewSelect().
		Model(&LockState{}).
		Where("infohash = ?", s).
		Exists(context.Background())

	if err != nil {
		DbL.Fatalf("failed to check lock state: %v", err)
	}

	return result
}
