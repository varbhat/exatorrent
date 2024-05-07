package db

import (
	"context"
	"github.com/uptrace/bun"
)

type Tracker struct {
	bun.BaseModel `bun:"table:trackerdb"`

	Url string `bun:"url,pk"`
}

type TrackerRepo struct {
	conn *bun.DB
}

func NewTrackerRepo(db *bun.DB) *TrackerRepo {
	return &TrackerRepo{conn: db}
}

func (tr *TrackerRepo) Open(dsn string) {
	_, err := tr.conn.
		NewCreateTable().
		Model(&Tracker{}).
		IfNotExists().
		Exec(context.Background())

	if err != nil {
		DbL.Printf("failed to create table tracker: %v", err)
	}
}

func (tr *TrackerRepo) Close() error {
	return tr.conn.Close()
}

func (tr *TrackerRepo) Add(s string) {
	_, err := tr.conn.
		NewInsert().
		Model(&Tracker{
			Url: s,
		}).
		On("conflict (url) do nothing").
		Exec(context.Background())

	if err != nil {
		DbL.Printf("failed to insert tracker: %v", err)
	}
}

func (tr *TrackerRepo) Delete(s string) {
	_, err := tr.conn.
		NewDelete().
		Model(&Tracker{}).
		Where("url = ?", s).
		Exec(context.Background())
	if err != nil {
		DbL.Printf("failed to delete tracker: %v", err)
	}
}

func (tr *TrackerRepo) DeleteN(count int) {

	subQuery := tr.conn.NewSelect().
		Model(&Tracker{}).
		Column("url").
		Limit(count)

	_, err := tr.conn.
		NewDelete().
		Model(&Tracker{}).
		Where("url in (?)", subQuery).
		Exec(context.Background())

	if err != nil {
		DbL.Printf("failed to delete tracker: %v", err)
	}
}

func (tr *TrackerRepo) DeleteAll() {
	_, err := tr.conn.
		NewDelete().
		Model(&Tracker{}).
		Where("1 = 1").
		Exec(context.Background())

	if err != nil {
		DbL.Printf("failed to delete tracker: %v", err)
	}
}

func (tr *TrackerRepo) Count() int64 {
	count, err := tr.conn.
		NewSelect().
		Model(&Tracker{}).
		Count(context.Background())

	if err != nil {
		DbL.Printf("failed to count tracker: %v", err)
	}

	return int64(count)
}

func (tr *TrackerRepo) Get() []string {
	var rs []string
	err := tr.conn.NewSelect().
		Model(&Tracker{}).
		Column("url").
		Scan(context.Background(), &rs)

	if err != nil {
		DbL.Printf("failed to fetch trackers: %v", err)
	}

	return rs
}
