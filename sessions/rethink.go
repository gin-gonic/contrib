package sessions

import (
	"github.com/boj/rethinkstore"
	"github.com/gorilla/sessions"
)

type RethinkStore interface {
	Store
}

// address: host:port
// db: database name
// table: table name
// maxIdle: maximum number of idle connections.
// maxOpen: maximum number of open connections.
// keyPairs: see https://godoc.org/github.com/gin-gonic/contrib/sessions#NewCookieStore
func NewRethinkStore(address, db, table string, maxIdle, maxOpen int, keyPairs []byte) (RethinkStore, error) {
	store, err := rethinkstore.NewRethinkStore(address, db, table, maxIdle, maxOpen, keyPairs)
	if err != nil {
		return nil, err
	}
	return &rethinkStore{store}, nil
}

type rethinkStore struct {
	*rethinkstore.RethinkStore
}

func (c *rethinkStore) Options(options Options) {
	c.RethinkStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
