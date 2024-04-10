package main

import (
	"log"

	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/routes"
	"nearbyassist/internal/server"
	"nearbyassist/internal/storage"
)

func main() {
	// Load configuration file
	config := config.LoadConfig()

	// Load file store
	store := storage.NewStorage(config)
	store.CreateDirectories()

	// Load database configuration
	db := db.NewDatabase(config)

	// Create and start the server
	server := server.NewServer(config, db, store)
	routes.RegisterRoutes(server)

	go server.Websocket.SaveMessages()
	go server.Websocket.ForwardMessages()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
