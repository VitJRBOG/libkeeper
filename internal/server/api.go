package server

import (
	"database/sql"
	"handbook-api/internal/db"
	"handbook-api/internal/models"
	"net/http"
	"net/url"
	"strconv"
)

func getNotes(dbConn *sql.DB, params url.Values) ([]models.Note, error) {
	var err error
	id := -1

	if params.Has("id") {
		id, err = strconv.Atoi(params.Get("id"))
		if err != nil {
			return nil, Error{http.StatusBadRequest, "'id' param must be integer"}
		}
	}

	notes, err := db.SelectNotes(dbConn, id)
	if err != nil {
		return nil, Error{http.StatusServiceUnavailable, "couldn't get notes"}
	}

	return notes, nil
}
