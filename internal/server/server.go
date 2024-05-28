package server

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/routing_engine"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/suggestion_engine"
	"nearbyassist/internal/utils"
	"nearbyassist/internal/websocket"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo             *echo.Echo
	Websocket        *websocket.Websocket
	DB               db.Database
	Storage          storage.Storage
	RouteEngine      routing_engine.Engine
	SuggestionEngine suggestion_engine.Engine
	Encrypt          encryption.Encryption
	Auth             authenticator.Authenticator
	Port             string
	AllowedOrigins   []string
}

func NewServer(conf *config.Config, ws *websocket.Websocket, db db.Database, storage storage.Storage, auth authenticator.Authenticator, router routing_engine.Engine, courtier suggestion_engine.Engine, crypto encryption.Encryption) *Server {
	NewServer := &Server{
		Echo:             echo.New(),
		Websocket:        ws,
		DB:               db,
		Storage:          storage,
		RouteEngine:      router,
		SuggestionEngine: courtier,
		Encrypt:          crypto,
		Auth:             auth,
		Port:             conf.Port,
		AllowedOrigins:   conf.AllowedOrigins,
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
