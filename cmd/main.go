package main

import (
	"log"

	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db/mysql"
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/hash"
	"nearbyassist/internal/id_generator"
	"nearbyassist/internal/routes"
	"nearbyassist/internal/routing_engine"
	"nearbyassist/internal/server"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/suggestion_engine"
	"nearbyassist/internal/websocket"
)

func main() {
	// Load configuration file
	config := config.LoadConfig()

	// Load file store
	store := storage.NewStorage(config)
	store.Initialize()

	// Load database configuration
	db := mysql.NewMysqlDatabase(config)

	serverConfig := server.ServerConfig{
		Config:           config,
		Websocket:        websocket.NewWebsocket(),
		DB:               db.Conn,
		Storage:          store,
		RouteEngine:      routing_engine.NewOSRM(config),
		SuggestionEngine: suggestion_engine.NewCourtier(),
		IdGen:            id_generator.NewNanoIdGenerator(),
		Encrypt:          encryption.NewAes(config),
		Hash:             hash.NewSha(),
		Auth:             authenticator.NewJWTAuthenticator(config),
	}

	// Create and start the server
	server, err := server.NewServer(serverConfig)
	if err != nil {
		log.Fatal(err)
	}

	routes.RegisterRoutes(server)

	go server.Websocket.SaveMessages()
	go server.Websocket.ForwardMessages()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
