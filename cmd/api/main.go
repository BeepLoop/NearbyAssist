package main

import (
	"log"

	"nearbyassist/internal/server"
)

func main() {

	server := server.NewServer()

	log.Println("starting server ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
