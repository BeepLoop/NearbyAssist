package server

import (
	"fmt"
	"nearbyassist/internal/handler"
	"net/http"
	"time"
)

var port = 8080

type Server struct {
	port int
}

func NewServer() *http.Server {

	NewServer := &Server{
		port: port,
	}

	handler := handler.Handler{}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(handler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
