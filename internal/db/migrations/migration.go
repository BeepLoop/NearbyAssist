package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	migrateDown := flag.Bool("down", false, "Specify true to run down migration")
	migrateUp := flag.Bool("up", false, "Specify true to run up migration")
	flag.Parse()

	godotenv.Load()
	var dsn string
	if os.Getenv("GO_ENV") == "development" {
		dsn = os.Getenv("MIGRATE_DSN_DEV")
	} else {
		dsn = os.Getenv("MIGRATE_DSN_PROD")
	}

	m, err := migrate.New("file:///home/johnloydmulit/go/personal/nearbyassist/internal/db/migrations/", dsn)
	if err != nil {
		panic("Error occurred while creating migration: " + err.Error())
	}

	if *migrateDown {
		err := m.Down()
		if err != nil {
			panic("Error occurred while migrating down: " + err.Error())
		}
		return
	}

	if *migrateUp {
		err := m.Up()
		if err != nil {
			panic("Error occurred while migrating up: " + err.Error())
		}
		return
	}

	fmt.Println("Please specify down or up migration by setting -down / -up flag to true")
}
