package db

import (
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	MAX_RETRIES = 10
)

func NewMysql(config mysql.Config) (*sqlx.DB, error) {
	retries := 0

	for {
		conn, err := sqlx.Connect("mysql", config.FormatDSN())
		if err != nil {
			log.Printf("Error connecting to database: %v\n", err)
			if retries >= MAX_RETRIES {
				return nil, err
			}

			log.Printf("Retrying in 5 seconds...\n")
			time.Sleep(5 * time.Second)
			continue
		}

		return conn, nil
	}
}
