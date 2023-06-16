package db

import (
	"database/sql"
	"fmt"
	"libkeeper-api/internal/models"

	_ "github.com/lib/pq" // Postgres driver
)

// Connection stores DB connection.
type Connection struct {
	Conn *sql.DB
}

// NewConnection creates a connection to the PostgreSQL database and returns the struct with it.
func NewConnection(dsn string) (Connection, error) {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return Connection{}, fmt.Errorf("unable to create a database connection: %s", err.Error())
	}

	err = dbConn.Ping()
	if err != nil {
		return Connection{}, fmt.Errorf("unable connect to database: %s", err.Error())
	}

	return Connection{
		Conn: dbConn,
	}, nil
}

// CreateNote inserts new entries into the "note" and "version" tables.
func CreateNote(dbConn Connection, note models.Note, version models.Version) error {
	query := "WITH new_note AS (INSERT INTO note(title, c_date) VALUES($1, $2) RETURNING id)" +
		"INSERT INTO version(full_text, c_date, checksum, note_id) VALUES(" +
		"$3, $4, $5, (SELECT id FROM new_note))"

	_, err := dbConn.Conn.Exec(query, note.Title, note.CreationDate, version.FullText,
		version.CreationDate, version.Checksum)
	if err != nil {
		return fmt.Errorf("failed to insert entries into the 'note' and 'version' tables: %s", err)
	}

	return nil
}
