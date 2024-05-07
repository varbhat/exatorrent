package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/uptrace/bun"
	"time"
)

type TorrentEntity struct {
	bun.BaseModel `bun:"table:torrent"`

	InfoHash  string    `bun:"infohash,pk"`
	Started   bool      `bun:"started"`
	AddedAt   time.Time `bun:"addedat,type:timestamptz"`
	StartedAt time.Time `bun:"startedat,type:timestamptz"`
}

type TorrentRepo struct {
	conn *bun.DB
}

func NewTorrentRepo(db *bun.DB) *TorrentRepo {
	return &TorrentRepo{conn: db}
}

func (tr *TorrentRepo) Open(s string) {
	var _, err = tr.conn.
		NewCreateTable().
		Model(&TorrentEntity{}).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		DbL.Fatalf("failed to create table: %v", err)
	}
}

func (tr *TorrentRepo) Close() error {
	return tr.conn.Close()
}

func (tr *TorrentRepo) Exists(hash metainfo.Hash) bool {
	result, err := tr.conn.
		NewSelect().
		Model(&TorrentEntity{}).
		Where("infohash = ?", hash.HexString()).
		Exists(context.Background())
	if err != nil {
		DbL.Printf("failed to check existence: %v", err)
	}
	return result
}

func (tr *TorrentRepo) Add(hash metainfo.Hash) error {
	_, err := tr.conn.
		NewInsert().
		Model(&TorrentEntity{
			InfoHash:  hash.HexString(),
			Started:   false,
			AddedAt:   time.Now(),
			StartedAt: time.Now(),
		}).
		On("conflict (infohash) do nothing").
		Exec(context.Background())
	return err
}

func (tr *TorrentRepo) Delete(hash metainfo.Hash) error {
	_, err := tr.conn.
		NewDelete().
		Model(&TorrentEntity{}).
		Where("infohash = ?", hash.HexString()).
		Exec(context.Background())
	return err
}

func (tr *TorrentRepo) Start(hash metainfo.Hash) error {
	_, err := tr.conn.
		NewUpdate().
		Model(&TorrentEntity{}).
		Where("infohash = ?", hash.HexString()).
		Set("started = ?", true).
		Set("startedat = ?", time.Now()).
		Exec(context.Background())
	return err
}

func (tr *TorrentRepo) SetStarted(hash metainfo.Hash, started bool) error {
	_, err := tr.conn.
		NewUpdate().
		Model(&TorrentEntity{}).
		Where("infohash = ?", hash.HexString()).
		Set("started = ?", started).
		Exec(context.Background())
	return err
}

func (tr *TorrentRepo) HasStarted(hash string) bool {
	var result bool
	var err = tr.conn.
		NewSelect().
		Model(&TorrentEntity{}).
		Column("started").
		Where("infohash = ?", hash).
		Scan(context.Background(), &result)

	if err == nil {
		return result
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false
	} else {
		DbL.Printf("failed to query started torrent: %v", err)
	}

	return false
}

func (tr *TorrentRepo) GetTorrent(hash metainfo.Hash) (*Torrent, error) {
	var te TorrentEntity
	err := tr.conn.NewSelect().
		Model(&te).
		Where("infohash = ?", hash.HexString()).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	var torrent = &Torrent{
		Infohash:  hash,
		Started:   te.Started,
		AddedAt:   te.AddedAt,
		StartedAt: te.StartedAt,
	}

	return torrent, nil
}

func (tr *TorrentRepo) GetTorrents() ([]*Torrent, error) {
	var teList = make([]*TorrentEntity, 0)
	err := tr.conn.NewSelect().
		Model(&teList).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	var rs = make([]*Torrent, 0, len(teList))
	for _, te := range teList {
		var ih metainfo.Hash
		ih, err = MetafromHex(te.InfoHash)
		if err != nil {
			return rs, err
		}
		rs = append(rs, &Torrent{
			Infohash:  ih,
			Started:   te.Started,
			AddedAt:   te.AddedAt,
			StartedAt: te.StartedAt,
		})
	}

	return rs, nil
}
