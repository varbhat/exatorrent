package db

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/uptrace/bun"
)

type FileState struct {
	bun.BaseModel `bun:"table:filestatedb"`

	FilePath string `bun:"filepath,unique:fs_uniq_key"`
	InfoHash string `bun:"infohash,unique:fs_uniq_key"`
}

type FileStateRepo struct {
	conn *bun.DB
}

func NewFileStateRepo(db *bun.DB) *FileStateRepo {
	return &FileStateRepo{conn: db}
}

func (fsr *FileStateRepo) Open(dsn string) {
	var _, err = fsr.conn.
		NewCreateTable().
		Model(&FileState{}).
		IfNotExists().
		Exec(context.Background())

	if err != nil {
		DbL.Fatalf("fail to create table: %v", err)
	}
}

func (fsr *FileStateRepo) Close() error {
	// do nothing close bun.db somewhere
	return fsr.conn.Close()
}

func (fsr *FileStateRepo) Add(fp string, hash metainfo.Hash) error {
	_, err := fsr.conn.NewInsert().
		Model(&FileState{
			FilePath: fp,
			InfoHash: hash.HexString(),
		}).
		On("conflict (filepath, infohash) do nothing").
		Exec(context.Background())
	return err
}

func (fsr *FileStateRepo) Get(hash metainfo.Hash) []string {
	var rs = make([]string, 0)
	err := fsr.conn.NewSelect().
		Model(&FileState{}).
		Column("filepath").
		Where("infohash = ?", hash.HexString()).
		Scan(context.Background(), &rs)

	if err != nil {
		DbL.Printf("fail to query with hash = %s. err : %v", hash.HexString(), err)
	}

	return rs
}

func (fsr *FileStateRepo) Deletefile(s string, hash metainfo.Hash) error {
	_, err := fsr.conn.NewDelete().
		Model(&FileState{}).
		Where("filepath = ?", s).
		Where("infohash = ?", hash.HexString()).
		Exec(context.Background())
	return err
}

func (fsr *FileStateRepo) Delete(hash metainfo.Hash) error {
	_, err := fsr.conn.NewDelete().
		Model(&FileState{}).
		Where("infohash = ?", hash.HexString()).
		Exec(context.Background())
	return err
}
