package admin

import (
	"nearbyassist/internal/service/auth"

	"github.com/jmoiron/sqlx"
)

type adminHandler struct {
	db      *sqlx.DB
	hash    auth.Hash
	encrypt auth.Encryption
}

func NewHandler(db *sqlx.DB, hash auth.Hash, encrypt auth.Encryption) *adminHandler {
	return &adminHandler{db: db, hash: hash, encrypt: encrypt}
}
