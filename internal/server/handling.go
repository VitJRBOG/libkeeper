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

// handler обрабатывает получаемые запросы.
func handler(dbConn *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sendData(w, map[string]map[string]string{"response": {"msg": "Hello world"}})
	})
}

func sendData(w http.ResponseWriter, values interface{}) {
	data, err := json.Marshal(values)
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

func sendError(w http.ResponseWriter, err error) {
	values := map[string]map[string]any{
		"error": {
			"status": http.StatusInternalServerError,
			"detail": "",
		},
	}

	if errInfo, ok := err.(Error); ok {
		values["error"]["status"] = errInfo.HTTPStatus
		values["error"]["detail"] = errInfo.Detail
	} else {
		values["error"]["detail"] = "internal server error"
	}

	data, err := json.Marshal(values)
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
