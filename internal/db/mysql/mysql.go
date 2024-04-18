package mysql

import (
	"fmt"
	"log"
	"nearbyassist/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	Conn *sqlx.DB
}

func NewMysqlDatabase(conf *config.Config) *Mysql {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.DB_User, conf.DB_Password, conf.DB_Host, conf.DB_Port, conf.DB_Name)

	instance := &Mysql{}

	for {
		conn, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			log.Printf("error connecting to database: %v\n", err)
			log.Printf("Retrying in 5 seconds...\n")
			time.Sleep(5 * time.Second)
			continue
		}

		instance.Conn = conn
		break
	}

	log.Printf("Connected to database %s\n", conf.DB_Name)
	return instance
}

func NewMysqlWithDb(db *sqlx.DB) *Mysql {
	return &Mysql{
		Conn: db,
	}
}
