package user

import "github.com/jmoiron/sqlx"

type MysqlUserStore struct {
	db *sqlx.DB
}

func NewMysqlUserStore(db *sqlx.DB) *MysqlUserStore {
	return &MysqlUserStore{
		db: db,
	}
}
