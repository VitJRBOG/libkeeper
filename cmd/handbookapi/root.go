package handbookapi

import (
	"fmt"
	"handbook-api/internal/config"
	"handbook-api/internal/db"
	"handbook-api/internal/server"
)

// Execute starts the main functions of program.
func Execute() {
	dbConnCfg := config.NewDBConnectionConfig()

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		dbConnCfg.DBMS(),
		dbConnCfg.DBMSUserName(), dbConnCfg.DBMSUserPassword(),
		dbConnCfg.Host(), dbConnCfg.Port(), dbConnCfg.DBName())

	dbConn := db.NewConnection(dsn)

	serverCfg := config.NewServerConfig()

	server.Up(dbConn, serverCfg)
}
