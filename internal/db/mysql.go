package db

import (
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	instance *sqlx.DB
)

func NewMysql(config mysql.Config) (*sqlx.DB, error) {
	if instance != nil {
		return instance, nil
	}

	conn, err := sqlx.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}
	instance = conn

	return instance, nil
}
