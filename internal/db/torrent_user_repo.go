package db

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/uptrace/bun"
)

type TorrentUser struct {
	bun.BaseModel `bun:"table:torrentuserdb"`

	Username string `bun:"username,unique:tu_uniq_key"`
	InfoHash string `bun:"infohash,unique:tu_uniq_key"`
}

type TorrentUserRepo struct {
	conn *bun.DB
}

func (tur *TorrentUserRepo) Open(s string) {
	var _, err = tur.conn.
		NewCreateTable().
		Model(&TorrentUser{}).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		DbL.Fatalf("failed to create table: %v", err)
	}
}

func (tur *TorrentUserRepo) Close() {
}

func (tur *TorrentUserRepo) Add(username string, hash metainfo.Hash) error {
	_, err := tur.conn.
		NewInsert().
		Model(&TorrentUser{
			Username: username,
			InfoHash: hash.HexString(),
		}).
		On("conflict (username, infohash) do nothing").
		Exec(context.Background())
	return err
}

func (tur *TorrentUserRepo) Remove(username string, hash metainfo.Hash) error {
	_, err := tur.conn.
		NewDelete().
		Model(&TorrentUser{}).
		Where("username = ?", username).
		Where("infohash = ?", hash.HexString()).
		Exec(context.Background())
	return err
}

func (tur *TorrentUserRepo) RemoveAll(username string) error {
	_, err := tur.conn.
		NewDelete().
		Model(&TorrentUser{}).
		Where("username = ?", username).
		Exec(context.Background())
	return err
}

func (tur *TorrentUserRepo) RemoveAllMi(hash metainfo.Hash) error {
	_, err := tur.conn.
		NewDelete().
		Model(&TorrentUser{}).
		Where("infohash = ?", hash.HexString()).
		Exec(context.Background())
	return err
}

func (tur *TorrentUserRepo) HasUser(username string, hash string) bool {
	result, err := tur.conn.
		NewSelect().
		Model(&TorrentUser{}).
		Where("username = ?", username).
		Where("infohash = ?", hash).
		Exists(context.Background())

	if err == nil {
		return result
	}

	DbL.Printf("failed to check if user exists: %v", err)
	return false
}

func (tur *TorrentUserRepo) ListTorrents(username string) []metainfo.Hash {
	var hashList []string
	err := tur.conn.
		NewSelect().
		Model(&TorrentUser{}).
		Column("infohash").
		Where("username = ?", username).
		Scan(context.Background(), &hashList)
	if err != nil {
		DbL.Printf("failed to list torrents: %v", err)
		return nil
	}
	var rs = make([]metainfo.Hash, 0, len(hashList))

	for _, hash := range hashList {
		ih, err := MetafromHex(hash)
		if err != nil {
			return rs
		}
		rs = append(rs, ih)
	}

	return rs
}

func (tur *TorrentUserRepo) ListUsers(hash metainfo.Hash) []string {
	var rs []string
	err := tur.conn.
		NewSelect().
		Model(&TorrentUser{}).
		Column("username").
		Where("infohash = ?", hash.HexString()).
		Scan(context.Background(), &rs)

	if err != nil {
		DbL.Printf("failed to list users: %v", err)
	}

	return rs
}

func NewTorrentUserRepo(db *bun.DB) *TorrentUserRepo {
	return &TorrentUserRepo{conn: db}
}
