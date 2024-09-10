package main

import (
	"log"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/routing_engine"
	"nearbyassist/internal/server"
	"nearbyassist/internal/service/auth"
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
		Encrypt:          auth.NewAES([]byte(config.EncryptionKey)),
		JWT:              auth.NewJWTAuthenticator(config.JwtSecret, config.JwtDuration),
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
