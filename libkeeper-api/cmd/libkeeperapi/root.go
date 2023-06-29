package libkeeperapi

import (
	"fmt"
	"libkeeper-api/internal/config"
	"libkeeper-api/internal/db"
	"libkeeper-api/internal/loggers"
	"libkeeper-api/internal/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	serverRepresentative := getCompletionHerald()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go osSignalsReception(&wg, serverRepresentative)

	wg.Add(1)
	go server.Up(&wg, serverRepresentative, serverCfg, dbConn)

	wg.Wait()
	loggers.NewInfoLogger().Println("program exited successfully")
}

func getCompletionHerald() chan os.Signal {
	serverRepresentative := make(chan os.Signal, 1)

	return serverRepresentative
}

func osSignalsReception(wg *sync.WaitGroup, serverRepresentative chan os.Signal) {
	signal.Notify(serverRepresentative, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg.Done()
}
