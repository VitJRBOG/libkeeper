package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// NewConnection create a new DB connection and return pointer on it.
func NewConnection(dsn string) *sql.DB {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("unable connect to database: %s\n", err.Error())
	}

	return dbConn
}
