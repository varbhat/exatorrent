package db

import (
	"math/rand"
	"testing"
)

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func genRandomString(n int) string {
	characterLength := len(characters)
	var buff = make([]byte, n)
	for i := 0; i < n; i++ {
		buff[i] = characters[rand.Int63n(int64(characterLength))]
	}
	return string(buff)
}

func TestTrackerRepo_Open(t *testing.T) {
	tr := NewTrackerRepo(db)

	tr.Open("")
}

func TestTrackerRepo_Add(t *testing.T) {
	tr := NewTrackerRepo(db)
	for i := 0; i < 16; i++ {
		tr.Add(genRandomString(32))
	}
}

func TestTrackerRepo_Count(t *testing.T) {
	tr := NewTrackerRepo(db)
	count := tr.Count()
	t.Log(count)
}

func TestTrackerRepo_Delete(t *testing.T) {
	tr := NewTrackerRepo(db)
	tr.Add("nSRTypJ9nqLTi5D6AeYd292Xn4zcactV")
	tr.Add("nSRTypJ9nqLTi5D6AeYd292Xn4zcactV")
	tr.Delete("nSRTypJ9nqLTi5D6AeYd292Xn4zcactV")
}

func TestTrackerRepo_DeleteAll(t *testing.T) {
	tr := NewTrackerRepo(db)
	tr.DeleteAll()
}

func TestTrackerRepo_DeleteN(t *testing.T) {
	tr := NewTrackerRepo(db)
	tr.DeleteN(4)
}

func TestTrackerRepo_Get(t *testing.T) {
	tr := NewTrackerRepo(db)
	var rs = tr.Get()
	t.Logf("%v", rs)
	t.Logf("count: %d", len(rs))
}
