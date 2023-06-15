package server

import (
	"fmt"
	"libkeeper-api/internal/config"
	"log"
	"net/http"
)

// Up starts the server.
func Up(serverCfg config.ServerCfg) {
	handling()
	log.Println("request handling is ready")

	address := fmt.Sprintf(":%s", serverCfg.Port)

	err := http.ListenAndServe(address, nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server launch error: %s", err)
	}
}
