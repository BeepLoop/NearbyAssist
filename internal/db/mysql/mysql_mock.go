package mysql

import (
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error occurred while opening a stub database connection: %v", err)
	}

	return db, mock
}
