package sessions

import (
	"testing"
)

const rethinkTestServer = "localhost:28015"

var newRethinkStore = func(_ *testing.T) Store {
	store, err := NewRethinkStore(rethinkTestServer, "testdb", "testtable", 5, 5, []byte("secret"))
	if err != nil {
		panic(err)
	}
	return store
}

func TestRethink_SessionGetSet(t *testing.T) {
	sessionGetSet(t, newRethinkStore)
}

func TestRethink_SessionDeleteKey(t *testing.T) {
	sessionDeleteKey(t, newRethinkStore)
}

func TestRethink_SessionFlashes(t *testing.T) {
	sessionFlashes(t, newRethinkStore)
}

func TestRethink_SessionClear(t *testing.T) {
	sessionClear(t, newRethinkStore)
}

func TestRethink_SessionOptions(t *testing.T) {
	sessionOptions(t, newRethinkStore)
}
