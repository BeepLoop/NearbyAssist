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

	// Load file disk
	disk := storage.NewDiskStorage(config)
	disk.Initialize()

	// Load database configuration
	mysql, err := db.NewMysql(mysql.Config{
		User:                 config.DB_USER,
		Passwd:               config.DB_PWD,
		Net:                  config.DB_NET,
		Addr:                 config.DB_HOST + ":" + config.DB_PORT,
		DBName:               config.DB_NAME,
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer mysql.Close()

	serverConfig := server.ServerConfig{
		Config:           config,
		DB:               mysql,
		Storage:          disk,
		JWT:              auth.NewJWTAuthenticator(config.JWT_SECRET, config.JWT_DURATION),
		Encrypt:          auth.NewAES([]byte(config.ENCRYPTION_KEY)),
		Websocket:        websocket.NewWebsocket(),
		RouteEngine:      routing_engine.NewOSRM(config),
		SuggestionEngine: suggestion_engine.NewCourtier(),
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
