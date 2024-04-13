package mysql

import (
	"nearbyassist/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCountVendor(t *testing.T) {
	sql, mock := newMock()
	sqlx := sqlx.NewDb(sql, "sqlmock")
	db := NewMysqlWithDb(sqlx)
	defer db.Conn.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	query := "SELECT COUNT\\(\\*\\) FROM Vendor WHERE restricted = 0"
	mock.ExpectQuery(query).WillReturnRows(rows)

	count, err := db.CountVendor(models.VENDOR_STATUS_UNRESTRICTED)

	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindVendorById(t *testing.T) {
	sql, mock := newMock()
	sqlx := sqlx.NewDb(sql, "sqlmock")
	db := NewMysqlWithDb(sqlx)
	defer db.Conn.Close()

	rows := sqlmock.NewRows([]string{"id", "vendorId", "rating", "job"}).
		AddRow(1, 1, "4", "Plumbing").
		AddRow(2, 2, "5", "Electrician")

	query := "SELECT id, vendorId, rating, job, restricted FROM Vendor WHERE id = ?"
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	vendor, err := db.FindVendorById(1)

	assert.NoError(t, err)
	assert.NotNil(t, vendor)
	assert.Equal(t, 1, vendor.Id)
	assert.Equal(t, "Plumbing", vendor.Job)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRestrictVendor(t *testing.T) {
	sql, mock := newMock()
	sqlx := sqlx.NewDb(sql, "sqlmock")
	db := NewMysqlWithDb(sqlx)
	defer db.Conn.Close()

	result := sqlmock.NewResult(0, 1)

	query := "UPDATE Vendor SET restricted = 1 WHERE vendorId = ?"
	mock.ExpectExec(query).WithArgs(1).WillReturnResult(result)

	err := db.RestrictVendor(1)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUnrestrictVendor(t *testing.T) {
	sql, mock := newMock()
	sqlx := sqlx.NewDb(sql, "sqlmock")
	db := NewMysqlWithDb(sqlx)
	defer db.Conn.Close()

	result := sqlmock.NewResult(0, 1)

	query := "UPDATE Vendor SET restricted = 0 WHERE vendorId = ?"
	mock.ExpectExec(query).WithArgs(1).WillReturnResult(result)

	err := db.UnrestrictVendor(1)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
