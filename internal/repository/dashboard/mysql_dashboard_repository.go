package dashboard_repo

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlDashboardRepository struct {
	db *sqlx.DB
}

func NewMysqlDashboardRepository(db *sqlx.DB) *MysqlDashboardRepository {
	return &MysqlDashboardRepository{db: db}
}

func (s *MysqlDashboardRepository) GetUserData() (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	countUserQuery := "SELECT COUNT(id) FROM User"
	userCount := 0
	if err := s.db.GetContext(ctx, &userCount, countUserQuery); err != nil {
		return nil, err
	}

	countVerifiedQuery := "SELECT COUNT(id) FROM IdentityVerification WHERE status = 'approved'"
	verifiedCount := 0
	if err := s.db.GetContext(ctx, &verifiedCount, countVerifiedQuery); err != nil {
		return nil, err
	}

	countExpertsQuery := "SELECT COUNT(vendorId) FROM Vendor"
	vendorCount := 0
	if err := s.db.GetContext(ctx, &vendorCount, countExpertsQuery); err != nil {
		return nil, err
	}

	countRestrictedQuery := "SELECT COUNT(userId) FROM Restricted"
	restrictedCount := 0
	if err := s.db.GetContext(ctx, &restrictedCount, countRestrictedQuery); err != nil {
		return nil, err
	}

	data := &models.UserData{
		Total:      userCount,
		Verified:   verifiedCount,
		Expert:     vendorCount,
		Restricted: restrictedCount,
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return data, nil
}

func (s *MysqlDashboardRepository) GetBugReportData() (*models.WeeklyBugReportData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reportCountThisWeekQuery := `
        WITH date_series AS (
            SELECT CURDATE() - INTERVAL n DAY AS reportDate
            FROM (
                SELECT 0 AS n UNION ALL
                SELECT 1 UNION ALL
                SELECT 2 UNION ALL
                SELECT 3 UNION ALL
                SELECT 4 UNION ALL
                SELECT 5 UNION ALL
                SELECT 6
            ) numbers
        )
        SELECT 
            d.reportDate AS date,
            COUNT(s.id) AS count
        FROM 
            date_series d
        LEFT JOIN 
            BugReport s
        ON 
            DATE(s.createdAt) = d.reportDate
        GROUP BY 
            d.reportDate
        ORDER BY 
            d.reportDate;
    `
	reportCountThisWeek := make([]models.DailyBugReportData, 0)
	if err := s.db.SelectContext(ctx, &reportCountThisWeek, reportCountThisWeekQuery); err != nil {
		return nil, err
	}

	reportCountLastWeekQuery := `
        WITH date_series AS (
            SELECT CURDATE() - INTERVAL (7 + n) DAY AS reportDate
            FROM (
                SELECT 0 AS n UNION ALL
                SELECT 1 UNION ALL
                SELECT 2 UNION ALL
                SELECT 3 UNION ALL
                SELECT 4 UNION ALL
                SELECT 5 UNION ALL
                SELECT 6
            ) numbers
        )
        SELECT 
            d.reportDate AS date,
            COUNT(s.id) AS count
        FROM 
            date_series d
        LEFT JOIN 
            BugReport s
        ON 
            DATE(s.createdAt) = d.reportDate
        GROUP BY 
            d.reportDate
        ORDER BY 
            d.reportDate;
    `
	reportCountLastWeek := make([]models.DailyBugReportData, 0)
	if err := s.db.SelectContext(ctx, &reportCountLastWeek, reportCountLastWeekQuery); err != nil {
		return nil, err
	}

	totalReportsThisWeek := 0
	for _, report := range reportCountThisWeek {
		totalReportsThisWeek += report.Count
	}

	totalReportsLastWeek := 0
	for _, report := range reportCountLastWeek {
		totalReportsLastWeek += report.Count
	}

	// NOTE:
	// Positive difference = bad
	// Negative difference = good
	difference := totalReportsThisWeek - totalReportsLastWeek

	data := &models.WeeklyBugReportData{
		Total:      totalReportsThisWeek,
		Daily:      reportCountThisWeek,
		Difference: difference,
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return data, nil
}

func (s *MysqlDashboardRepository) GetVendorReportData() (*models.WeeklyVendorReportData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reportCountThisWeekQuery := `
        WITH date_series AS (
            SELECT CURDATE() - INTERVAL n DAY AS reportDate
            FROM (
                SELECT 0 AS n UNION ALL
                SELECT 1 UNION ALL
                SELECT 2 UNION ALL
                SELECT 3 UNION ALL
                SELECT 4 UNION ALL
                SELECT 5 UNION ALL
                SELECT 6
            ) numbers
        )
        SELECT 
            d.reportDate AS date,
            COUNT(ru.id) AS count
        FROM 
            date_series d
        LEFT JOIN 
            UserReport ru
        ON 
            DATE(ru.createdAt) = d.reportDate
        GROUP BY 
            d.reportDate
        ORDER BY 
            d.reportDate;
    `
	reportCountThisWeek := make([]models.DailyVendorReportData, 0)
	if err := s.db.SelectContext(ctx, &reportCountThisWeek, reportCountThisWeekQuery); err != nil {
		return nil, err
	}

	reportCountLastWeekQuery := `
        WITH date_series AS (
            SELECT CURDATE() - INTERVAL (7 + n) DAY AS reportDate
            FROM (
                SELECT 0 AS n UNION ALL
                SELECT 1 UNION ALL
                SELECT 2 UNION ALL
                SELECT 3 UNION ALL
                SELECT 4 UNION ALL
                SELECT 5 UNION ALL
                SELECT 6
            ) numbers
        )
        SELECT 
            d.reportDate AS date,
            COUNT(ru.id) AS count
        FROM 
            date_series d
        LEFT JOIN 
            UserReport ru
        ON 
            DATE(ru.createdAt) = d.reportDate
        GROUP BY 
            d.reportDate
        ORDER BY 
            d.reportDate;
    `
	reportCountLastWeek := make([]models.DailyVendorReportData, 0)
	if err := s.db.SelectContext(ctx, &reportCountLastWeek, reportCountLastWeekQuery); err != nil {
		return nil, err
	}

	totalReportsThisWeek := 0
	for _, report := range reportCountThisWeek {
		totalReportsThisWeek += report.Count
	}

	totalReportsLastWeek := 0
	for _, report := range reportCountLastWeek {
		totalReportsLastWeek += report.Count
	}

	// NOTE:
	// Positive difference = bad
	// Negative difference = good
	difference := totalReportsThisWeek - totalReportsLastWeek

	data := &models.WeeklyVendorReportData{
		Total:      totalReportsThisWeek,
		Daily:      reportCountThisWeek,
		Difference: difference,
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return data, nil
}

func (s *MysqlDashboardRepository) GetIdentityVerificationRequestsData() (*models.IdentityVerificationRequestData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	countPendingQuery := `
        SELECT
            COUNT(id)
        FROM
            IdentityVerification
        WHERE
            status = 'pending'
    `
	pendingCount := 0
	if err := s.db.GetContext(ctx, &pendingCount, countPendingQuery); err != nil {
		return nil, err
	}

	data := &models.IdentityVerificationRequestData{
		Total: pendingCount,
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return data, nil
}

func (s *MysqlDashboardRepository) GetVendorApplicationRequestsData() (*models.VendorApplicationRequestData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	countPendingQuery := `
        SELECT
            Count(id)
        FROM
            Application
        WHERE
            status = 'pending'
    `
	pendingCount := 0
	if err := s.db.GetContext(ctx, &pendingCount, countPendingQuery); err != nil {
		return nil, err
	}

	data := &models.VendorApplicationRequestData{
		Total: pendingCount,
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return data, nil
}

// ==== NEW

func (s *MysqlDashboardRepository) TotalUsers() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM User"

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlDashboardRepository) TotalVendors() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(*) FROM Vendor"

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlDashboardRepository) TotalReported() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            COUNT(*)
        FROM
            User u
            LEFT JOIN UserReport r ON u.id = r.reportedUserId
        WHERE
            r.reportedUserId IS NULL
    `

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlDashboardRepository) TotalRestricted() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            COUNT(*)
        FROM
            User u
            LEFT JOIN Restricted r ON u.id = r.userId
        WHERE
            r.userId IS NOT NULL
    `

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlDashboardRepository) RecentUsers() ([]*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            id, name, email, imageUrl, createdAt
        FROM
            User
        ORDER BY
            createdAt DESC
        LIMIT 10
    `

	users := make([]*models.UserModel, 0)
	if err := s.db.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return users, nil
}

func (s *MysqlDashboardRepository) TotalServices() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(*) FROM Service"

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlDashboardRepository) TotalActiveServices() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(*) FROM Service WHERE disabled = 0"

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlDashboardRepository) TotalPendingApplications() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(*) FROM Application WHERE status = 'pending'"

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlDashboardRepository) TotalActiveReports() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT COUNT(*) FROM UserReport WHERE status = 'pending'"

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlDashboardRepository) GetBookingData() (*dto.WeeklyBookingData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookingCountThisWeekQuery := `
        WITH date_series AS (
            SELECT CURDATE() - INTERVAL n DAY AS reportDate
            FROM (
                SELECT 0 AS n UNION ALL
                SELECT 1 UNION ALL
                SELECT 2 UNION ALL
                SELECT 3 UNION ALL
                SELECT 4 UNION ALL
                SELECT 5 UNION ALL
                SELECT 6
            ) numbers
        )
        SELECT 
            d.reportDate AS date,
            COUNT(t.id) AS count
        FROM 
            date_series d
        LEFT JOIN 
            Booking t
        ON 
            DATE(t.createdAt) = d.reportDate
        GROUP BY 
            d.reportDate
        ORDER BY 
            d.reportDate;
    `

	bookingCountThisWeek := make([]dto.DailyBookingData, 0)
	if err := s.db.SelectContext(ctx, &bookingCountThisWeek, bookingCountThisWeekQuery); err != nil {
		return nil, err
	}

	bookingCountLastWeekQuery := `
        WITH date_series AS (
            SELECT CURDATE() - INTERVAL n DAY AS reportDate
            FROM (
                SELECT 0 AS n UNION ALL
                SELECT 1 UNION ALL
                SELECT 2 UNION ALL
                SELECT 3 UNION ALL
                SELECT 4 UNION ALL
                SELECT 5 UNION ALL
                SELECT 6
            ) numbers
        )
        SELECT 
            d.reportDate AS date,
            COUNT(t.id) AS count
        FROM 
            date_series d
        LEFT JOIN 
            Booking t
        ON 
            DATE(t.createdAt) = d.reportDate
        GROUP BY 
            d.reportDate
        ORDER BY 
            d.reportDate;
    `

	bookingCountLastWeek := make([]models.DailyBookingData, 0)
	if err := s.db.SelectContext(ctx, &bookingCountLastWeek, bookingCountLastWeekQuery); err != nil {
		return nil, err
	}

	totalThisWeek := 0
	for _, t := range bookingCountThisWeek {
		totalThisWeek += t.Count
	}

	totalLastWeek := 0
	for _, t := range bookingCountLastWeek {
		totalLastWeek += t.Count
	}

	// NOTE:
	// Positive difference = bad
	// Negative difference = good
	difference := totalThisWeek - totalLastWeek

	data := &dto.WeeklyBookingData{
		Total:      totalThisWeek,
		Daily:      bookingCountThisWeek,
		Difference: difference,
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return data, nil
}
