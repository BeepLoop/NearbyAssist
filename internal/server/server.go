package server

import (
	"nearbyassist/internal/websocket"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo      *echo.Echo
	Websocket *websocket.Websocket
}

func NewServer() *Server {

	ws := websocket.NewWebsocket()

	NewServer := &Server{
		Echo:      echo.New(),
		Websocket: ws,
	}

	return NewServer
}

func (s *Server) Start(listenAddr string) error {
	return s.Echo.Start(listenAddr)
}
