package server

import (
	"fmt"
	"net/http"
	"os"
	"pc-beantragung/internal/database"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {
	getPort := func() int {
		port, valid := os.LookupEnv("PORT")
		if valid {
			port, _ := strconv.Atoi(port)
			return port
		}
		return 8080
	}

	NewServer := &Server{
		port: getPort(),
		db:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
