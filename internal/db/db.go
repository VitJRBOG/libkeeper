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
