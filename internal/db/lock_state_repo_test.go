package db

import "testing"

func TestLockStateRepo_Open(t *testing.T) {
	fsr := NewLockStateRepo(db)

	fsr.Open("")
}

func TestLockStateRepo_Lock(t *testing.T) {
	fsr := NewLockStateRepo(db)

	err := fsr.Lock(ih1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLockStateRepo_Unlock(t *testing.T) {
	fsr := NewLockStateRepo(db)

	err := fsr.Unlock(ih1)
	if err != nil {
		t.Fatal(err)
	}

	locked := fsr.IsLocked(ih1.HexString())
	t.Logf("locked: %t", locked)
}

func TestLockStateRepo_IsLocked(t *testing.T) {
	fsr := NewLockStateRepo(db)

	locked := fsr.IsLocked(ih1.HexString())
	t.Logf("locked: %t", locked)
}
