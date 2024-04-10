package db

import (
	"log"
	"nearbyassist/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Conn *sqlx.DB
}

func NewDatabase(conf *config.Config) *DB {
	conn, err := sqlx.Connect("mysql", conf.DSN)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	return &DB{Conn: conn}
}
