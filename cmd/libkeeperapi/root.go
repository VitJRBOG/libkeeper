package libkeeperapi

import (
	"fmt"
	"libkeeper-api/internal/config"
	"libkeeper-api/internal/db"
	"libkeeper-api/internal/loggers"
	"libkeeper-api/internal/server"
	"log"
)

// Execute starts the main functions of program.
func Execute() {
	loggers.InitializeDefaultLogger()

	dbConnectionCfg, err := config.NewDBConnectionCfg()
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConnectionCfg.User, dbConnectionCfg.Password,
		dbConnectionCfg.HostAddress, dbConnectionCfg.HostPort,
		dbConnectionCfg.DBName,
		dbConnectionCfg.SSLMode)

	dbConn, err := db.NewConnection(dsn)
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	serverCfg, err := config.NewServerCfg()
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	server.Up(serverCfg, dbConn)
}
