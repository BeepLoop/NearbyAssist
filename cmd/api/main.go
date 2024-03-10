package main

import (
	"log"
	"os"

	"nearbyassist/internal/config"
	"nearbyassist/internal/db"
	"nearbyassist/internal/server"
)

func init() {
	if err := config.Init(); err != nil {
		log.Fatal("error reading environment variables: ", err)
	}

	if err := db.Init(); err != nil {
		log.Fatal("error initializing database connection: ", err)
	}

	if err := os.MkdirAll("store", 0777); err != nil {
		log.Fatal("Unable to initialize file store: ", err)
	}
}

func main() {

	server := server.NewServer()

	log.Println("starting server ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
