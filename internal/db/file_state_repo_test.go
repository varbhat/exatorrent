package db

import (
	"log"
	"testing"
)

func TestFileStateRepo_Open(t *testing.T) {
	fsr := NewFileStateRepo(db)

	fsr.Open("")
}

func TestFileStateRepo_Add(t *testing.T) {
	var err error
	fsr := NewFileStateRepo(db)
	err = fsr.Add("fp", ih1)
	if err != nil {
		log.Fatal(err)
	}

	// check duplicate
	err = fsr.Add("fp", ih1)
	if err != nil {
		log.Fatal(err)
	}

	err = fsr.Add("fp1", ih1)
	if err != nil {
		log.Fatal(err)
	}

	err = fsr.Add("fp2", ih2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileStateRepo_Get(t *testing.T) {
	fsr := NewFileStateRepo(db)

	rs := fsr.Get(ih1)
	t.Log(rs)
}

func TestFileStateRepo_Deletefile(t *testing.T) {
	fsr := NewFileStateRepo(db)

	err := fsr.Deletefile("fp", ih1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileStateRepo_Delete(t *testing.T) {

	fsr := NewFileStateRepo(db)

	err := fsr.Delete(ih1)
	if err != nil {
		t.Fatal(err)
	}
}
