package main

import (
	"log"

	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/hash"
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
	db := db.NewDatabase(config)

	// Load authenticator configuration
	auth := authenticator.NewJWTAuthenticator(config)

	// Load websocket configuration
	ws := websocket.NewWebsocket(db)

	// Load Routing Engine configuration
	engine := routing_engine.NewOSRM(config)

	// Load Suggestion Engine configuration
	courtier := suggestion_engine.NewCourtier()

	// Load encryption configuration
	crypto := encryption.NewAes(config)

	// Load hashing algorithm
	hash := hash.NewSha()

	// Create and start the server
	server := server.NewServer(config, ws, db, store, auth, engine, courtier, crypto, hash)
	routes.RegisterRoutes(server)

	go server.Websocket.SaveMessages()
	go server.Websocket.ForwardMessages()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
