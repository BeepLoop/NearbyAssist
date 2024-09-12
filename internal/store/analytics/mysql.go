package analytics

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlAnalyticsStore struct {
	db *sqlx.DB
}

func NewMysqlAnalyticsStore(db *sqlx.DB) *MysqlAnalyticsStore {
	return &MysqlAnalyticsStore{
		db: db,
	}
}

func (s *MysqlAnalyticsStore) UserCount(filter UserStatusFilter) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(id) FROM User"

	switch filter {
	case USER_STATUS_VERIFIED:
		query += " WHERE verified = 1"
	case USER_STATUS_UNVERIFIED:
		query += " WHERE verified = 0"
	case USER_STATUS_ALL:
	}

	count := 0
	err := s.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlAnalyticsStore) VendorCount(filter VendorStatusFilter) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(id) FROM Vendor"

	switch filter {
	case VENDOR_STATUS_RESTRICTED:
		query += " WHERE restricted = 1"
	case VENDOR_STATUS_UNRESTRICTED:
		query += " WHERE restricted = 0"
	case VENDOR_STATUS_ALL:
	}

	count := 0
	err := s.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlAnalyticsStore) ApplicationCount(filter ApplicationStatusFilter) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(id) FROM Application"

	switch filter {
	case APPLICATION_STATUS_PENDING:
		query += " WHERE status = 'pending'"
	case APPLICATION_STATUS_APPROVED:
		query += " WHERE status = 'approved'"
	case APPLICATION_STATUS_REJECTED:
		query += " WHERE status = 'rejected'"
	}

	count := 0
	err := s.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlAnalyticsStore) ComplaintCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(id) FROM Complaint"

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}
