package server

import (
	"fmt"
	"libkeeper-api/internal/config"
	"libkeeper-api/internal/db"
	"libkeeper-api/internal/loggers"
	"log"
	"net"
	"net/http"
	"time"
)

// Up starts the server.
func Up(serverCfg config.ServerCfg, dbConn db.Connection) {
	infoLogger := loggers.NewInfoLogger()

	handling(dbConn)
	infoLogger.Println("request handling is ready")

	address := fmt.Sprintf(":%s", serverCfg.Port)

	err := http.ListenAndServe(address, logging(http.DefaultServeMux, infoLogger))
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server launch error: %s", err)
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
