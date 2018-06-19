package state

import (
	"github.com/go-pg/pg"
)

// OpenConnection creates a PostgreSQL connection with configured connection parameters
func OpenConnection(user, password, database, host string) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: database,
	})

	return db, nil
}
