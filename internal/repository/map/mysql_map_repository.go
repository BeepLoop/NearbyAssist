package map_repo

import "github.com/jmoiron/sqlx"

type MysqlMapRepository struct {
	db *sqlx.DB
}

func NewMysqlMapRepository(db *sqlx.DB) *MysqlMapRepository {
	return &MysqlMapRepository{db: db}
}
