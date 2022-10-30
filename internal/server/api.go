package server

import (
	"database/sql"
	"fmt"
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

func updateNote(dbConn *sql.DB, params url.Values) (int, error) {
	var err error
	noteID := -1

	if params.Has("id") {
		noteID, err = strconv.Atoi(params.Get("id"))
		if err != nil {
			return -1, Error{http.StatusBadRequest, "'id' param must be integer"}
		}
	} else {
		return -1, Error{http.StatusBadRequest, "'id' param is empty"}
	}

	if !params.Has("title") {
		return -1, Error{http.StatusBadRequest, "'title' param is empty"}
	}

	if !params.Has("text") {
		return -1, Error{http.StatusBadRequest, "'text' param is empty"}
	}

	if !params.Has("date") {
		return -1, Error{http.StatusBadRequest, "'date' param is empty"}
	}

	if !params.Has("checksum") {
		return -1, Error{http.StatusBadRequest, "'checksum' param is empty"}
	}

	date, err := time.Parse("2006-01-02 15:04:05", params.Get("date"))
	if err != nil {
		return -1, Error{http.StatusBadRequest, "'date' has the invalid format"}
	}

	note := models.Note{
		ID:    noteID,
		Title: params.Get("title"),
	}

	changesNumber, err := db.UpdateNote(dbConn, note)
	if err != nil {
		return -1, Error{http.StatusServiceUnavailable, "couldn't update note"}
	}

	if changesNumber == 0 {
		return -1, Error{http.StatusBadRequest, fmt.Sprintf(
			"no rows with id = %d were found", noteID)}
	}

	version := models.Version{
		Text:     params.Get("text"),
		Date:     date.Unix(),
		Checksum: params.Get("checksum"),
		NoteID:   noteID,
	}

	versionID, err := db.InsertVersion(dbConn, version)
	if err != nil {
		return -1, Error{http.StatusServiceUnavailable, "couldn't create version"}
	}

	return versionID, nil
}

func deleteNote(dbConn *sql.DB, params url.Values) error {
	var err error
	noteID := -1

	if params.Has("id") {
		noteID, err = strconv.Atoi(params.Get("id"))
		if err != nil {
			return Error{http.StatusBadRequest, "'id' param must be integer"}
		}
	} else {
		return Error{http.StatusBadRequest, "'id' param is empty"}
	}

	version := models.Version{
		NoteID: noteID,
	}

	changesNumber, err := db.DeleteVersionsByNoteID(dbConn, version)

	if changesNumber == 0 {
		return Error{http.StatusBadRequest,
			fmt.Sprintf("no rows with note_id = %d were found", noteID)}
	}

	note := models.Note{
		ID: noteID,
	}

	changesNumber, err = db.DeleteNote(dbConn, note)
	if err != nil {
		return Error{http.StatusServiceUnavailable, "couldn't update note"}
	}

	if changesNumber == 0 {
		return Error{http.StatusBadRequest,
			fmt.Sprintf("no rows with id = %d were found", noteID)}
	}

	return nil
}