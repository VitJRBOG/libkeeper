package server

import (
	"database/sql"
	"log"
	"net/http"
)

// handler обрабатывает получаемые запросы.
func handler(dbConn *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			log.Printf("error: %s: %s", r.URL.String(), err.Error())
		}
	})
}
