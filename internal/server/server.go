package server

import (
	"nearbyassist/internal/config"
	"nearbyassist/internal/routing_engine"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/suggestion_engine"
	"nearbyassist/internal/utils"
	"nearbyassist/internal/websocket"
	"os"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type ServerConfig struct {
	Config           *config.Config
	Websocket        *websocket.Websocket
	DB               *sqlx.DB
	Storage          storage.Storage
	RouteEngine      routing_engine.Engine
	SuggestionEngine suggestion_engine.Engine
	Encrypt          auth.Encryption
	JWT              auth.Authenticator
}

type Server struct {
	LOG_FILE *os.File

	Echo           *echo.Echo
	Port           string
	AllowedOrigins []string

	Websocket        *websocket.Websocket
	DB               *sqlx.DB
	Storage          storage.Storage
	RouteEngine      routing_engine.Engine
	SuggestionEngine suggestion_engine.Engine
	Encrypt          auth.Encryption
	JWT              auth.Authenticator
}

func NewServer(options ServerConfig) (*Server, error) {
	file, err := os.OpenFile(options.Config.LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	NewServer := &Server{
		LOG_FILE: file,

		Echo:           echo.New(),
		Port:           options.Config.Port,
		AllowedOrigins: options.Config.AllowedOrigins,

		Websocket:        options.Websocket,
		DB:               options.DB,
		Storage:          options.Storage,
		RouteEngine:      options.RouteEngine,
		SuggestionEngine: options.SuggestionEngine,
		Encrypt:          options.Encrypt,
		JWT:              options.JWT,
	}

	return NewServer, nil
}

func (s *Server) Start() error {
	s.Echo.Validator = &utils.Validator{Validator: validator.New()}

	s.middlewares()
	s.routes()

	if err := s.Echo.Start(":" + s.Port); err != nil {
		s.LOG_FILE.Close()
		return err
	}

	return nil
}
