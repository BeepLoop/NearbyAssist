package server

import (
	"nearbyassist/internal/config"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/email"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/service/route_engine"
	"nearbyassist/internal/service/suggestion_engine"
	"nearbyassist/internal/service/websocket"

	"github.com/jmoiron/sqlx"
)

type ServerConfig struct {
	Config *config.Config

	WS *websocket.Websocket

	Mailman email.MailService

	DB *sqlx.DB
	FS fs.FileStorage

	RouteEngine      route_engine.Engine
	SuggestionEngine suggestion_engine.Engine

	Hash    auth.Hash
	Encrypt auth.Encryption
	JWT     auth.Authenticator
}
