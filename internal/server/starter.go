package server

import (
	"database/sql"
	"fmt"
	"handbook-api/internal/config"
	"log"
	"net/http"
	"time"
)

// Up поднимает сервер.
func Up(dbConn *sql.DB, serverCfg config.ServerConfig) {
	handler(dbConn)

	addr := fmt.Sprintf(":%s", serverCfg.Port())

	err := http.ListenAndServe(addr, logging(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("server error: %s\n", err.Error())
	}
}

// logging пишет обрабатываемые запросы в лог.
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begins := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(begins)

		log.Printf("[%s] %s %s", r.Method, r.RequestURI, timeElapsed)
	})
}
