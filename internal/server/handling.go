package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Error store info about error.
type Error struct {
	HTTPStatus int
	Detail     string
}

// Error return a text representation of error info.
func (e Error) Error() string {
	return fmt.Sprintf("status %d: %s", e.HTTPStatus, e.Detail)
}

// handler handles received requests.
func handler(dbConn *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sendData(w, http.StatusOK, "Hello world")
	})

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				log.Println(err.Error())
				sendError(w, err)
				return
			}

			noteID, versionID, err := createNote(dbConn, r.Form)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, http.StatusCreated, []map[string]int{{
				"note_id":    noteID,
				"version_id": versionID,
			}})

		case http.MethodGet:
			err := r.ParseForm()
			if err != nil {
				log.Println(err.Error())
				sendError(w, err)
				return
			}

			notes, err := getNotes(dbConn, r.Form)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, http.StatusOK, notes)
		case http.MethodPut:
			log.Println("update note")
			// TODO: describe a PUT method handler for /notes
		case http.MethodDelete:
			log.Println("delete note")
			// TODO: describe a DELETE method handler for /notes
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})
}

func sendData(w http.ResponseWriter, status int, values interface{}) {
	response := map[string]interface{}{
		"status":   status,
		"response": values,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}
}

func sendError(w http.ResponseWriter, reqError error) {
	response := map[string]interface{}{
		"status": http.StatusInternalServerError,
		"error":  "internal server error",
	}

	if errInfo, ok := reqError.(Error); ok {
		response["status"] = errInfo.HTTPStatus
		response["error"] = errInfo.Detail
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
