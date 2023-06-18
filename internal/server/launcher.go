package server

import (
	"context"
	"fmt"
	"libkeeper-api/internal/config"
	"libkeeper-api/internal/db"
	"libkeeper-api/internal/loggers"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// Up starts the server.
func Up(wg *sync.WaitGroup, signalToExit chan os.Signal,
	serverCfg config.ServerCfg, dbConn db.Connection) {
	infoLogger := loggers.NewInfoLogger()
	srv := serverSettingUp(serverCfg, infoLogger)

	go waitForExitSignal(signalToExit, srv, dbConn, infoLogger)

	handling(dbConn)
	infoLogger.Println("request handling is ready")

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server launch error: %s", err)
	}

	wg.Done()
}

func serverSettingUp(serverCfg config.ServerCfg, infoLogger *log.Logger) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverCfg.Port),
		Handler: logging(http.DefaultServeMux, infoLogger),
	}

	return srv
}

func waitForExitSignal(signalToExit chan os.Signal, srv *http.Server,
	dbConn db.Connection, infoLogger *log.Logger) {
	<-signalToExit

	serverShuttingDown(srv, infoLogger)
	closeDBConnection(dbConn)
}

func serverShuttingDown(srv *http.Server, infoLogger *log.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server shutdown failed: %s", err)
	}

	infoLogger.Println("server exited successfully")
}

func closeDBConnection(dbConn db.Connection) {
	err := dbConn.Conn.Close()
	if err != nil {
		log.Printf("error when closing the database connection: %s", err)
	}
}

func logging(next http.Handler, infoLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begins := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(begins)

		ip, err := getIP(r)
		if err != nil {
			log.Println(err)
		}
		infoLogger.Printf("[%s] [%s] %s %s", ip, r.Method, r.RequestURI, timeElapsed)
	})
}

func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("failed to receive ip: %s", err)
	}
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	return "no ip", fmt.Errorf("no valid IP found")
}
