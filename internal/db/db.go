package db

import (
	"nearbyassist/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Connection *sqlx.DB

func Init() error {
	conn, err := sqlx.Connect("mysql", config.Env.DSN)
	if err != nil {
		return err
	}
	Connection = conn

	return nil
}
