package config

import (
	"log"
	"os"
)

// DBConnection store a DB connection params.
type DBConnectionConfig struct {
	dbms             string
	host             string
	port             string
	dbmsUserName     string
	dbmsUserPassword string
	dbName           string
}

// NewDBConnectionConfig parse the environment variables and return a DBConnectionConfig struct.
func NewDBConnectionConfig() DBConnectionConfig {
	dbms := os.Getenv("DBMS")
	host := os.Getenv("DBMS_HOST")
	port := os.Getenv("DBMS_PORT")
	dbmsUserName := os.Getenv("DBMS_USER_NAME")
	dbmsUserPassword := os.Getenv("DMBS_USER_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	someIsEmpty := false

	if dbms == "" {
		log.Println("DBMS env variable is empty")
		someIsEmpty = true
	}
	if host == "" {
		log.Println("DBMS_HOST env variable is empty")
		someIsEmpty = true
	}
	if port == "" {
		log.Println("DBMS_PORT env variable is empty")
		someIsEmpty = true
	}
	if dbmsUserName == "" {
		log.Println("DBMS_USER_NAME env variable is empty")
		someIsEmpty = true
	}
	if dbmsUserPassword == "" {
		log.Println("DMBS_USER_PASSWORD env variable is empty")
		someIsEmpty = true
	}
	if dbName == "" {
		log.Println("DB_NAME env is empty")
		someIsEmpty = true
	}

	if someIsEmpty {
		log.Fatalln("some desktop environments is empty")
	}

	return DBConnectionConfig{
		dbms:             dbms,
		host:             host,
		port:             port,
		dbmsUserName:     dbmsUserName,
		dbmsUserPassword: dbmsUserPassword,
		dbName:           dbName,
	}
}

// DBMS return the dbms field of DBConnectionConfig.
func (c *DBConnectionConfig) DBMS() string {
	return c.dbms
}

// Host return the host field of DBConnectionConfig.
func (c *DBConnectionConfig) Host() string {
	return c.host
}

// Port return the port field of DBConnectionConfig.
func (c *DBConnectionConfig) Port() string {
	return c.port
}

// DBMSUserName return the dbmsUserName field of DBConnectionConfig.
func (c *DBConnectionConfig) DBMSUserName() string {
	return c.dbmsUserName
}

// DBMSUserPassword return the dbmsUserPassword field of DBConnectionConfig.
func (c *DBConnectionConfig) DBMSUserPassword() string {
	return c.dbmsUserPassword
}

// DBName return the dbName field of DBConnectionConfig.
func (c *DBConnectionConfig) DBName() string {
	return c.dbName
}

// ServerConfig store a server params.
type ServerConfig struct {
	port string
}

// NewServerConfig parse the environment variables and return a ServerConfig struct.
func NewServerConfig() ServerConfig {
	port := os.Getenv("SERVER_PORT")

	if port == "" {
		log.Fatalln("SERVER_PORT env variable is empty")
	}

	return ServerConfig{
		port: port,
	}
}

// Port return the port field of ServerConfig.
func (c *ServerConfig) Port() string {
	return c.port
}
