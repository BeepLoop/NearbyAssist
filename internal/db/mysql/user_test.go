package mysql

import (
	"database/sql"
	"log"
	"nearbyassist/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error occurred while opening a stub database connection: %v", err)
	}

	return db, mock
}

func TestFindUserById(t *testing.T) {
	u := &models.UserModel{
		Model: models.Model{
			Id: 1,
		},
		Name:     "john loyd mulit",
		Email:    "jlmulit68@gmail.com",
		ImageUrl: "https://example.com",
	}

	db, mock := newMock()
	sqlxDb := sqlx.NewDb(db, "sqlmock")
	d := NewMysqlWithDb(sqlxDb)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "imageUrl"}).
		AddRow(u.Id, u.Name, u.Email, u.ImageUrl)

	query := "SELECT id, name, email, imageUrl FROM User WHERE id = ?"
	mock.ExpectQuery(query).WithArgs(u.Id).WillReturnRows(rows)

	user, err := d.FindUserById(u.Id)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindUserByIdError(t *testing.T) {
	u := &models.UserModel{
		Model: models.Model{
			Id: 1,
		},
		Name:     "john loyd mulit",
		Email:    "jlmulit68@gmail.com",
		ImageUrl: "https://example.com",
	}

	db, mock := newMock()
	sqlxDb := sqlx.NewDb(db, "sqlmock")
	d := NewMysqlWithDb(sqlxDb)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "imageUrl"})

	query := "SELECT id, name, email, imageUrl FROM User WHERE id = ?"
	mock.ExpectQuery(query).WithArgs(u.Id).WillReturnRows(rows)

	user, err := d.FindUserById(u.Id)

	assert.Nil(t, user)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindUserByEmail(t *testing.T) {
	u := &models.UserModel{
		Model: models.Model{
			Id: 1,
		},
		Name:     "john loyd mulit",
		Email:    "jlmulit68@gmail.com",
		ImageUrl: "https://example.com",
	}

	db, mock := newMock()
	sqlxDb := sqlx.NewDb(db, "sqlmock")
	d := NewMysqlWithDb(sqlxDb)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "imageUrl"}).
		AddRow(u.Id, u.Name, u.Email, u.ImageUrl)

	query := "SELECT id, name, email, imageUrl FROM User WHERE email = ?"
	mock.ExpectQuery(query).WithArgs(u.Email).WillReturnRows(rows)

	user, err := d.FindUserByEmail(u.Email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john loyd mulit", user.Name)
	assert.NotEqual(t, 0, user.Id)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindUserByEmailError(t *testing.T) {
	u := &models.UserModel{
		Model: models.Model{
			Id: 1,
		},
		Name:     "john loyd mulit",
		Email:    "jlmulit68@gmail.com",
		ImageUrl: "https://example.com",
	}

	db, mock := newMock()
	sqlxDb := sqlx.NewDb(db, "sqlmock")
	d := NewMysqlWithDb(sqlxDb)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "imageUrl"})

	query := "SELECT id, name, email, imageUrl FROM User WHERE email = ?"
	mock.ExpectQuery(query).WithArgs(u.Email).WillReturnRows(rows)

	user, err := d.FindUserByEmail(u.Email)

	assert.Error(t, err)
	assert.Nil(t, user)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCountUserShouldBeZero(t *testing.T) {
	db, mock := newMock()
	sqlxDb := sqlx.NewDb(db, "sqlmock")
	d := NewMysqlWithDb(sqlxDb)

	rows := sqlmock.NewRows([]string{"count"})

	query := "SELECT COUNT\\(\\*\\) FROM User"
	mock.ExpectQuery(query).WillReturnRows(rows)

	count, err := d.CountUser()

	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestCountUserShouldBeOne(t *testing.T) {
	db, mock := newMock()
	sqlxDb := sqlx.NewDb(db, "sqlmock")
	d := NewMysqlWithDb(sqlxDb)

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(1)

	query := "SELECT COUNT\\(\\*\\) FROM User"
	mock.ExpectQuery(query).WillReturnRows(rows)

	count, err := d.CountUser()

	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
