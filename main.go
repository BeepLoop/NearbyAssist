package main

import (
	"log"

	"nearbyassist/internal/config"
	"nearbyassist/internal/db/mysql"
	"nearbyassist/internal/routes"
	"nearbyassist/internal/server"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/websocket"
)

func main() {
	// Load configuration file
	config := config.LoadConfig()

	// Load file store
	store := storage.NewDiskStorage(config)
	store.CreateDirectories()

	// Load database configuration
	db := mysql.NewMysqlDatabase(config)

	ws := websocket.NewWebsocket(db)

	// Create and start the server
	server := server.NewServer(config, ws, db, store)
	routes.RegisterRoutes(server)

	go server.Websocket.SaveMessages()
	go server.Websocket.ForwardMessages()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
