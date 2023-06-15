package server

import (
	"fmt"
	"libkeeper-api/internal/config"
	"libkeeper-api/internal/db"
	"log"
	"net/http"
)

// Up starts the server.
func Up(serverCfg config.ServerCfg, dbConn db.Connection) {
	handling(dbConn)
	log.Println("request handling is ready")

	address := fmt.Sprintf(":%s", serverCfg.Port)

	err := http.ListenAndServe(address, nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server launch error: %s", err)
	}
}
