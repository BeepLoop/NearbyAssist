package main

import (
	"log"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/routing_engine"
	"nearbyassist/internal/server"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/suggestion_engine"
	"nearbyassist/internal/websocket"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration file
	config := config.LoadConfig()

	// Load hashing algorithm
	hash := auth.NewSha256()

	// Load file disk
	storage := fs.NewDiskStorage(config, hash)

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
		Config: config,

		DB: mysql,
		FS: storage,

		Websocket: websocket.NewWebsocket(),

		JWT:     auth.NewJWTAuthenticator(config.JWT_SECRET, config.JWT_DURATION),
		Encrypt: auth.NewAES([]byte(config.ENCRYPTION_KEY)),
		Hash:    hash,

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
