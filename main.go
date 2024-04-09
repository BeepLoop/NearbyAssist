package main

import (
	"log"
	"os"

	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/routes"
	"nearbyassist/internal/server"
)

func init() {
	if err := config.Init(); err != nil {
		log.Fatal("error reading environment variables: ", err)
	}

	if err := db.Init(); err != nil {
		log.Fatal("error initializing database connection: ", err)
	}

	if err := os.MkdirAll("store/application", 0777); err != nil {
		log.Fatal("Unable to initialize file store: ", err)
	}

	if err := os.MkdirAll("store/service", 0777); err != nil {
		log.Fatal("Unable to initialize file store: ", err)
	}
}

func main() {
	server := server.NewServer()

	routes.RegisterRoutes(server)

	go server.Websocket.SaveMessages()
	go server.Websocket.ForwardMessages()

	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
