package main

import (
	"log"

	// "nearbyassist/internal/authenticator"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/encryption"
	// "nearbyassist/internal/hash"
	// "nearbyassist/internal/id_generator"
	"nearbyassist/internal/routing_engine"
	"nearbyassist/internal/server"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/suggestion_engine"
	"nearbyassist/internal/websocket"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration file
	config := config.LoadConfig()

	// Load file store
	store := storage.NewStorage(config)
	store.Initialize()

	// Load database configuration
	mysql, err := db.NewMysql(mysql.Config{
		User:                 config.DB_User,
		Passwd:               config.DB_Password,
		Net:                  "tcp",
		Addr:                 config.DB_Host + ":" + config.DB_Port,
		DBName:               config.DB_Name,
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer mysql.Close()

	serverConfig := server.ServerConfig{
		Config:           config,
		Websocket:        websocket.NewWebsocket(),
		DB:               mysql,
		Storage:          store,
		RouteEngine:      routing_engine.NewOSRM(config),
		SuggestionEngine: suggestion_engine.NewCourtier(),
		// IdGen:            id_generator.NewNanoIdGenerator(),
		Encrypt: encryption.NewAes(config),
		// Hash:             hash.NewSha(),
		// Auth:             authenticator.NewJWTAuthenticator(config),
	}

	// Create and start the server
	server, err := server.NewServer(serverConfig)
	if err != nil {
		log.Fatal(err)
	}

	go server.Websocket.SaveMessages()
	go server.Websocket.ForwardMessages()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
