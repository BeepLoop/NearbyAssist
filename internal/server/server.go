package server

import (
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/utils"
	"nearbyassist/internal/websocket"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo      *echo.Echo
	Websocket *websocket.Websocket
	DB        db.Database
	Storage   storage.Storage
	Port      string
}

func NewServer(conf *config.Config, ws *websocket.Websocket, db db.Database, storage storage.Storage) *Server {
	NewServer := &Server{
		Echo:      echo.New(),
		Websocket: ws,
		DB:        db,
		Storage:   storage,
		Port:      conf.Port,
	}

	return NewServer
}

func (s *Server) configure() {
	s.Echo.Validator = &utils.Validator{Validator: validator.New()}
}

func (s *Server) Start() error {
	s.configure()
	s.registerMiddleware()

	return s.Echo.Start(":" + s.Port)
}
