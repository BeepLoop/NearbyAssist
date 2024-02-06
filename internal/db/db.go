package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB_CONN *sqlx.DB

func Init() error {
	conn, err := sqlx.Connect("mysql", "root:Password_1@/nearby_assist")
	if err != nil {
		return err
	}
	DB_CONN = conn

	return nil
}
