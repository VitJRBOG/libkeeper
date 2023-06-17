package server

import (
	"encoding/json"
	"fmt"
	"libkeeper-api/internal/db"
	"libkeeper-api/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Error stores info about error.
type Error struct {
	HTTPStatus int
	Detail     string
}

// Error returns a text representation of error info.
func (e Error) Error() string {
	return fmt.Sprintf("status %d: %s", e.HTTPStatus, e.Detail)
}

func handling(dbConn db.Connection) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			log.Printf("something went wrong: %s", err)
		}
	})

	http.HandleFunc("/note", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			err := parseRequestParams(r)
			if err != nil {
				sendError(w, err)
				return
			}

			versions, err := getVersions(dbConn, r)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, http.StatusOK, versions)
		case http.MethodPost:
			err := parseRequestParams(r)
			if err != nil {
				sendError(w, err)
				return
			}

			err = createNote(dbConn, r)
			if err != nil {
				sendError(w, err)
				return
			}
		case http.MethodPut:
			err := parseRequestParams(r)
			if err != nil {
				sendError(w, err)
				return
			}

			err = updateNote(dbConn, r)
			if err != nil {
				sendError(w, err)
				return
			}
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			notes, err := getNotes(dbConn)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, http.StatusOK, notes)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})
}

func sendData(w http.ResponseWriter, status int, values interface{}) {
	response := map[string]interface{}{
		"response": values,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}
}

func sendError(w http.ResponseWriter, reqError error) {
	response := map[string]interface{}{}

	if errInfo, ok := reqError.(Error); ok {
		w.WriteHeader(errInfo.HTTPStatus)
		response["error"] = errInfo.Detail
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		response["error"] = "internal server error"
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func parseRequestParams(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		log.Printf("unable to parse request params: %s", err)
		return Error{
			http.StatusInternalServerError,
			"unable to parse request params",
		}
	}

	return nil
}

func createNote(dbConn db.Connection, r *http.Request) error {
	emptyParam := []string{}
	someIsEmpty := false

	if r.PostFormValue("c_date") == "" {
		someIsEmpty = true
		emptyParam = append(emptyParam, "c_date")
	}

	if r.PostFormValue("checksum") == "" {
		someIsEmpty = true
		emptyParam = append(emptyParam, "checksum")
	}

	if someIsEmpty {
		errMsg := "some request parameters are missing: "

		for i := range emptyParam {
			if i > 0 && i < len(emptyParam) {
				errMsg += ", "
			}
			errMsg += fmt.Sprintf("'%s'", emptyParam[i])
		}

		return Error{
			http.StatusBadRequest,
			errMsg,
		}
	}

	creationDate, err := time.Parse("2006-01-02 15:04:05 -0700", r.PostFormValue("c_date"))
	if err != nil {
		return Error{
			http.StatusBadRequest,
			"the 'c_date' parameter must be in the format 'yyyy-mm-dd hh:mm:ss -0000'",
		}
	}
	cDate := strconv.FormatInt(creationDate.Unix(), 10)

	note := models.Note{
		Title:        r.PostFormValue("title"),
		CreationDate: cDate,
	}

	version := models.Version{
		FullText:     r.PostFormValue("full_text"),
		CreationDate: cDate,
		Checksum:     r.PostFormValue("checksum"),
	}

	err = db.CreateNote(dbConn, note, version)
	if err != nil {
		log.Printf("failed to create a new note: %s", err)
		return Error{
			http.StatusInternalServerError,
			"failed to create a new note",
		}
	}

	return nil
}

func getNotes(dbConn db.Connection) ([]models.Note, error) {
	notes, err := db.SelectNotes(dbConn)
	if err != nil {
		log.Printf("unable fetch notes from the database: %s", err)
		return nil, Error{
			http.StatusInternalServerError,
			"unable fetch notes from the database",
		}
	}

	return notes, nil
}

func getVersions(dbConn db.Connection, r *http.Request) ([]models.Version, error) {
	if r.FormValue("note_id") == "" {
		return nil, Error{
			http.StatusBadRequest,
			"the 'note_id' parameter is empty",
		}
	}

	noteID, err := strconv.Atoi(r.FormValue("note_id"))
	if err != nil {
		return nil, Error{
			http.StatusBadRequest,
			"the 'note_id' parameter must be an integer",
		}
	}

	versions, err := db.SelectVersions(dbConn, noteID)
	if err != nil {
		log.Printf("unable to fetch note versions from the database: %s", err)
		return nil, Error{
			http.StatusInternalServerError,
			"unable to fetch note versions from the database",
		}
	}

	return versions, nil
}

func updateNote(dbConn db.Connection, r *http.Request) error {
	emptyParam := []string{}
	someIsEmpty := false

	if r.PostFormValue("note_id") == "" {
		someIsEmpty = true
		emptyParam = append(emptyParam, "note_id")
	}

	if r.PostFormValue("c_date") == "" {
		someIsEmpty = true
		emptyParam = append(emptyParam, "c_date")
	}

	if r.PostFormValue("checksum") == "" {
		someIsEmpty = true
		emptyParam = append(emptyParam, "checksum")
	}

	if someIsEmpty {
		errMsg := "some request parameters are missing: "

		for i := range emptyParam {
			if i > 0 && i < len(emptyParam) {
				errMsg += ", "
			}
			errMsg += fmt.Sprintf("'%s'", emptyParam[i])
		}

		return Error{
			http.StatusBadRequest,
			errMsg,
		}
	}

	noteID, err := strconv.Atoi(r.PostFormValue("note_id"))
	if err != nil {
		return Error{
			http.StatusBadRequest,
			"the 'note_id' parameter must be an integer",
		}
	}

	creationDate, err := time.Parse("2006-01-02 15:04:05 -0700", r.PostFormValue("c_date"))
	if err != nil {
		return Error{
			http.StatusBadRequest,
			"the 'c_date' parameter must be in the format 'yyyy-mm-dd hh:mm:ss -0000'",
		}
	}
	cDate := strconv.FormatInt(creationDate.Unix(), 10)

	note := models.Note{
		ID:    noteID,
		Title: r.PostFormValue("title"),
	}

	version := models.Version{
		FullText:     r.PostFormValue("full_text"),
		CreationDate: cDate,
		Checksum:     r.PostFormValue("checksum"),
		NoteID:       noteID,
	}

	err = db.UpdateNote(dbConn, note, version)
	if err != nil {
		log.Printf("failed to update a note: %s", err)
		return Error{
			http.StatusInternalServerError,
			"failed to update a note",
		}
	}

	return nil
}