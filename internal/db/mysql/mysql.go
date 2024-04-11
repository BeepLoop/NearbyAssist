package mysql

import (
	"log"
	"nearbyassist/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	Conn *sqlx.DB
}

func NewMysqlDatabase(conf *config.Config) *Mysql {
	conn, err := sqlx.Connect("mysql", conf.DSN)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	return &Mysql{Conn: conn}
}
