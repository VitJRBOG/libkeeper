package libkeeperapi

import (
	"libkeeper-api/internal/config"
	"libkeeper-api/internal/server"
	"log"
)

// Execute starts the main functions of program.
func Execute() {
	serverCfg, err := config.NewServerCfg()
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	server.Up(serverCfg)
}
