package server

import (
	"nearbyassist/internal/utils"
	"nearbyassist/internal/websocket"

	"github.com/go-playground/validator"
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

func (s *Server) configure() {
	s.Echo.Validator = &utils.Validator{Validator: validator.New()}
}

func (s *Server) Start(listenAddr string) error {
	s.configure()
	s.registerMiddleware()

	return s.Echo.Start(listenAddr)
}
