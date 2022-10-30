package server

import (
	"database/sql"
	"handbook-api/internal/db"
	"handbook-api/internal/models"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func createNote(dbConn *sql.DB, params url.Values) (int, int, error) {
	if !params.Has("title") {
		return -1, -1, Error{http.StatusBadRequest, "'title' param is empty"}
	}

	if !params.Has("text") {
		return -1, -1, Error{http.StatusBadRequest, "'text' param is empty"}
	}

	if !params.Has("date") {
		return -1, -1, Error{http.StatusBadRequest, "'date' param is empty"}
	}

	if !params.Has("checksum") {
		return -1, -1, Error{http.StatusBadRequest, "'checksum' param is empty"}
	}

	date, err := time.Parse("2006-01-02 15:04:05", params.Get("date"))
	if err != nil {
		return -1, -1, Error{http.StatusBadRequest, "'date' has the invalid format"}
	}

	note := models.Note{
		Title: params.Get("title"),
		Date:  date.Unix(),
	}

	noteID, err := db.InsertNote(dbConn, note)
	if err != nil {
		return -1, -1, Error{http.StatusServiceUnavailable, "couldn't create note"}
	}

	version := models.Version{
		Text:     params.Get("text"),
		Date:     date.Unix(),
		Checksum: params.Get("checksum"),
		NoteID:   noteID,
	}

	versionID, err := db.InsertVersion(dbConn, version)
	if err != nil {
		return -1, -1, Error{http.StatusServiceUnavailable, "couldn't create version"}
	}

	return noteID, versionID, nil
}

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
