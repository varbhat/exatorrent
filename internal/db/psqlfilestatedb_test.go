package db

import (
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"log"
	"testing"
)

const pgDsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

var pgFileStateDb = &PsqlFsDb{}

func TestMain(m *testing.M) {
	var err error
	postgres := embeddedpostgres.NewDatabase()
	err = postgres.Start()
	if err != nil {
		log.Fatalf("failed to start postgres: %v", err)
	}

	defer func() {
		err = postgres.Stop()
		if err != nil {
			log.Fatalf("failed to stop postgres: %v", err)
		}
	}()

	pgFileStateDb.Open(pgDsn)
	code := m.Run()
	log.Printf("end with code %d", code)
}

func TestPsqlFsDb_Add(t *testing.T) {
	var err = pgFileStateDb.Add("fp", ih1)
	if err != nil {
		t.Fatalf("failed to add file state: %v", err)
	}
}
