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

// SelectNotes selects entries from the "note" table.
func SelectNotes(dbConn Connection) ([]models.Note, error) {
	query := "SELECT * FROM note"

	rows, err := dbConn.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to select entries from the 'note' table: %s", err)
	}

	notes := []models.Note{}

	for rows.Next() {
		note := models.Note{}

		err := rows.Scan(&note.ID, &note.Title, &note.CreationDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'note' table: %s", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

// SelectVersions selects entries from the "version" table by the value of "note_id" field. Returns them sorted by DESC.
func SelectVersions(dbConn Connection, noteID int) ([]models.Version, error) {
	query := "SELECT * FROM version WHERE note_id = $1 ORDER BY c_date DESC"

	rows, err := dbConn.Conn.Query(query, noteID)
	if err != nil {
		return nil, fmt.Errorf("failed to select entries from the 'version' table: %s", err)
	}

	versions := []models.Version{}

	for rows.Next() {
		version := models.Version{}

		err := rows.Scan(&version.ID, &version.FullText, &version.CreationDate, &version.Checksum, &version.NoteID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows from the 'version' table: %s", err)
		}

		versions = append(versions, version)
	}

	return versions, nil
}

// UpdateNote updates an existing entry in the 'note' table and inserts a new entry in the 'version' table.
func UpdateNote(dbConn Connection, note models.Note, version models.Version) error {
	query := "WITH updated_note AS (UPDATE note SET title=$1 WHERE id=$2 RETURNING id) " +
		"INSERT INTO version(full_text, c_date, checksum, note_id) " +
		"SELECT * FROM (SELECT $3 AS full_text, $4 AS c_date, $5 AS checksum, " +
		"(SELECT id FROM updated_note) AS note_id) AS new_version " +
		"WHERE NOT EXISTS (SELECT id FROM version WHERE version.note_id = (SELECT id FROM updated_note) " +
		"AND version.checksum = new_version.checksum)"

	_, err := dbConn.Conn.Exec(query, note.Title, note.ID, version.FullText, version.CreationDate, version.Checksum)
	if err != nil {
		return fmt.Errorf("failed to update the 'note' table entry and insert a new 'version' table entry: %s", err)
	}

	return nil
}

// DeleteNote deletes an existing entry from the 'note' table and deletes the associated entries from the 'version' table.
func DeleteNote(dbConn Connection, noteID int) error {
	query := "WITH deleted_versions AS (DELETE FROM version WHERE note_id = $1 RETURNING note_id) " +
		"DELETE FROM note WHERE id IN (SELECT note_id FROM deleted_versions)"

	_, err := dbConn.Conn.Exec(query, noteID)
	if err != nil {
		return fmt.Errorf("failed to delete the 'note' table entry: %s", err)
	}

	return nil
}
