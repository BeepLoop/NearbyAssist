package server

import (
	"encoding/gob"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/service/mailer"
	"nearbyassist/internal/service/route_engine"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"os"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	LOG_FILE *os.File

	Domain string

	Echo           *echo.Echo
	Port           string
	AllowedOrigins []string

	WS *websocket.Websocket

	DB *sqlx.DB
	FS fs.FileStorage

	Mailer mailer.Mailer

	RouteEngine      route_engine.Engine
	SuggestionEngine suggestion_engine.Engine

	Hash    core.Hash
	Encrypt core.Encryption
	JWT     core.Authenticator
}

func NewServer(options ServerConfig) (*Server, error) {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(options.Config.LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	NewServer := &Server{
		Domain: options.Config.DOMAIN,

		Echo:           echo.New(),
		Port:           options.Config.PORT,
		AllowedOrigins: options.Config.ALLOWED_ORIGINS,
		LOG_FILE:       file,

		WS: options.WS,

		DB: options.DB,
		FS: options.FS,

		Mailer: options.Mailer,

		Hash:    options.Hash,
		Encrypt: options.Encrypt,
		JWT:     options.JWT,

		RouteEngine:      options.RouteEngine,
		SuggestionEngine: options.SuggestionEngine,
	}

	return NewServer, nil
}

func (s *Server) Start() error {
	gob.Register(models.AdminModel{})

	s.Echo.Validator = &utils.Validator{Validator: validator.New()}

	s.middlewares()
	s.routes()

	s.WS.Start()

	if err := s.Echo.Start(":" + s.Port); err != nil {
		s.LOG_FILE.Close()
		s.WS.Stop()
		return err
	}

	return nil
}
