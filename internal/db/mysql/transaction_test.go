package mysql

import (
	"nearbyassist/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCountTransaction(t *testing.T) {
	tests := []struct {
		status        models.TransactionStatus
		query         string
		count         int
		expectedCount int
	}{
		{
			status:        models.TRANSACTION_STATUS_ONGOING,
			query:         "SELECT COUNT\\(\\*\\) FROM Transaction WHERE status = ?",
			count:         1,
			expectedCount: 1,
		},
		{
			status:        models.TRANSACTION_STATUS_CANCELLED,
			query:         "SELECT COUNT\\(\\*\\) FROM Transaction WHERE status = ?",
			count:         10,
			expectedCount: 10,
		},
	}

	for _, test := range tests {
		sql, mock := newMock()
		sqlx := sqlx.NewDb(sql, "sqlmock")
		db := NewMysqlWithDb(sqlx)
		defer db.Conn.Close()

		rows := sqlmock.NewRows([]string{"count"}).AddRow(test.count)

		mock.ExpectQuery(test.query).WillReturnRows(rows)

		count, err := db.CountTransaction(test.status)

		assert.NoError(t, err)
		assert.Equal(t, test.count, count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}
}
