package server

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/config"
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/hash"
	"nearbyassist/internal/id_generator"
	"nearbyassist/internal/routing_engine"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/suggestion_engine"
	"nearbyassist/internal/utils"
	"nearbyassist/internal/websocket"

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
	IdGen            id_generator.IdGenerator
	Encrypt          encryption.Encryption
	Hash             hash.Hash
	Auth             authenticator.Authenticator
}

type Server struct {
	Echo             *echo.Echo
	Websocket        *websocket.Websocket
	DB               *sqlx.DB
	Storage          storage.Storage
	RouteEngine      routing_engine.Engine
	SuggestionEngine suggestion_engine.Engine
	IdGen            id_generator.IdGenerator
	Encrypt          encryption.Encryption
	Hash             hash.Hash
	Auth             authenticator.Authenticator
	Port             string
	AllowedOrigins   []string
}

func NewServer(options ServerConfig) *Server {
	NewServer := &Server{
		Echo:             echo.New(),
		Websocket:        options.Websocket,
		DB:               options.DB,
		Storage:          options.Storage,
		RouteEngine:      options.RouteEngine,
		SuggestionEngine: options.SuggestionEngine,
		IdGen:            options.IdGen,
		Encrypt:          options.Encrypt,
		Hash:             options.Hash,
		Auth:             options.Auth,
		Port:             options.Config.Port,
		AllowedOrigins:   options.Config.AllowedOrigins,
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
