package handbookapi

import (
	"fmt"
	"handbook-api/internal/config"
	"handbook-api/internal/db"
	"handbook-api/internal/server"
	"log"
)

// Execute starts the main functions of program.
func Execute() {
	initializeLogger()

	dbConnCfg := config.NewDBConnectionConfig()

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		dbConnCfg.DBMS(),
		dbConnCfg.DBMSUserName(), dbConnCfg.DBMSUserPassword(),
		dbConnCfg.Host(), dbConnCfg.Port(), dbConnCfg.DBName(),
		dbConnCfg.SSLMode())

	dbConn := db.NewConnection(dsn)

	serverCfg := config.NewServerConfig()

	server.Up(dbConn, serverCfg)
}

func initializeLogger() {
	log.SetFlags(log.Ldate | log.Llongfile)
}
