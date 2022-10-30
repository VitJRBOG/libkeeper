package db

import (
	"database/sql"
	"handbook-api/internal/models"
	"log"

	_ "github.com/lib/pq" // Postgres driver
)

// NewConnection create a new DB connection and return pointer on it.
func NewConnection(dsn string) *sql.DB {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("unable connect to database: %s\n", err.Error())
	}

	return dbConn
}

// InsertNote inserts new row to 'notes' table and return ID of this new.
func InsertNote(dbConn *sql.DB, note models.Note) (int, error) {
	query := "INSERT INTO notes(title, c_date) VALUES($1, $2) RETURNING id"

	id := -1

	err := dbConn.QueryRow(query, note.Title, note.Date).Scan(&id)
	if err != nil {
		log.Printf("%s: %s", err.Error(), query)
		return -1, err
	}

	return id, nil
}

// SelectNotes selects field from 'notes' table by 'id', if id != -1.
// Else selects all fields from 'notes'.
func SelectNotes(dbConn *sql.DB, id int) ([]models.Note, error) {
	notes := []models.Note{}

	query := "SELECT * FROM notes"
	params := []any{}

	if id != -1 {
		query += " WHERE id = $1"
		params = append(params, id)
	}

	rows, err := dbConn.Query(query, params...)
	if err != nil {
		log.Printf("%s: %s", err.Error(), query)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf(err.Error())
		}
	}()

	for rows.Next() {
		note := models.Note{}

		if err := rows.Scan(&note.ID, &note.Title, &note.Date); err != nil {
			log.Printf(err.Error())
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

// InsertVersion inserts new row to 'versions' table and return ID of this new.
func InsertVersion(dbConn *sql.DB, version models.Version) (int, error) {
	query := "INSERT INTO versions(text, c_date, note_id) VALUES($1, $2, $3) RETURNING id"

	id := -1

	err := dbConn.QueryRow(query, version.Text, version.Date, version.NoteID).Scan(&id)
	if err != nil {
		log.Printf("%s: %s", err.Error(), query)
		return -1, err
	}

	return id, nil
}
