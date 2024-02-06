package main

import (
	"log"

	"nearbyassist/internal/db"
	"nearbyassist/internal/server"
)

func init() {
	if err := db.Init(); err != nil {
		log.Fatal("error initializing database connection: ", err)
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
