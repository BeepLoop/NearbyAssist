package db

import (
	"nearbyassist/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB_CONN *sqlx.DB

func Init() error {
	conn, err := sqlx.Connect("mysql", config.Env.DSN)
	if err != nil {
		return err
	}
	DB_CONN = conn

	return nil
}
