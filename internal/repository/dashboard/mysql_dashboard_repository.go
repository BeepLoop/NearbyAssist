package dashboard_repo

import (
	"context"
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

	countVerifiedQuery := "SELECT COUNT(id) FROM User WHERE verified = 1"
	verifiedCount := 0
	if err := s.db.GetContext(ctx, &verifiedCount, countVerifiedQuery); err != nil {
		return nil, err
	}

	countExpertsQuery := "SELECT COUNT(id) FROM Vendor"
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
            SystemComplaint s
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
            SystemComplaint s
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
            COUNT(v.id) AS count
        FROM 
            date_series d
        LEFT JOIN 
            VendorComplaint v
        ON 
            DATE(v.createdAt) = d.reportDate
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
            COUNT(v.id) AS count
        FROM 
            date_series d
        LEFT JOIN 
            VendorComplaint v
        ON 
            DATE(v.createdAt) = d.reportDate
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

func (s *MysqlDashboardRepository) GetTransactionData() (*models.WeeklyTransactionData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	transactionCountThisWeekQuery := `
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
            Transaction t
        ON 
            DATE(t.createdAt) = d.reportDate
        GROUP BY 
            d.reportDate
        ORDER BY 
            d.reportDate;
    `

	transactionCountThisWeek := make([]models.DailyTransactionData, 0)
	if err := s.db.SelectContext(ctx, &transactionCountThisWeek, transactionCountThisWeekQuery); err != nil {
		return nil, err
	}

	transactionCountLastWeekQuery := `
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
            Transaction t
        ON 
            DATE(t.createdAt) = d.reportDate
        GROUP BY 
            d.reportDate
        ORDER BY 
            d.reportDate;
    `

	transactionCountLastWeek := make([]models.DailyTransactionData, 0)
	if err := s.db.SelectContext(ctx, &transactionCountLastWeek, transactionCountLastWeekQuery); err != nil {
		return nil, err
	}

	totalThisWeek := 0
	for _, t := range transactionCountThisWeek {
		totalThisWeek += t.Count
	}

	totalLastWeek := 0
	for _, t := range transactionCountLastWeek {
		totalLastWeek += t.Count
	}

	// NOTE:
	// Positive difference = bad
	// Negative difference = good
	difference := totalThisWeek - totalLastWeek

	data := &models.WeeklyTransactionData{
		Total:      totalThisWeek,
		Daily:      transactionCountThisWeek,
		Difference: difference,
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return data, nil
}
