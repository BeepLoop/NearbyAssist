package mysql

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"time"
)

func (m *Mysql) CountApplication(status models.ApplicationStatus) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM Application"

	switch status {
	case models.APPLICATION_STATUS_PENDING:
		query += " WHERE status = 'pending'"
	case models.APPLICATION_STATUS_APPROVED:
		query += " WHERE status = 'approved'"
	case models.APPLICATION_STATUS_REJECTED:
		query += " WHERE status = 'rejected'"
	}

	count := 0
	err := m.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (m *Mysql) CreateApplication(application *request.NewApplication) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO
            Application (applicantId, job, latitude, longitude)
        VALUES
            (:applicantId, :job, :latitude, :longitude)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, application)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) FindApplicationById(id int) (*models.ApplicationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id, applicantId, job, status, latitude, longitude
        FROM
            Application
        WHERE
            id = ?
    `

	application := models.NewApplicationModel()
	err := m.Conn.GetContext(ctx, application, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return application, nil
}

func (m *Mysql) FindAllApplication(status models.ApplicationStatus) ([]response.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, applicantId, status, createdAt FROM Application"

	switch status {
	case models.APPLICATION_STATUS_PENDING:
		query += " WHERE status = 'pending'"
	case models.APPLICATION_STATUS_APPROVED:
		query += " WHERE status = 'approved'"
	case models.APPLICATION_STATUS_REJECTED:
		query += " WHERE status = 'rejected'"
	}

	applications := make([]response.Application, 0)
	err := m.Conn.SelectContext(ctx, &applications, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return applications, nil
}

func (m *Mysql) ApproveApplication(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := m.Conn.Beginx()
	if err != nil {
		return err
	}

	updateStatus := "UPDATE Application SET status = 'approved' WHERE id = ?"

	if _, err := tx.ExecContext(ctx, updateStatus, id); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.New("Failed to approve application and rollback transaction")
		}

		return err
	}

	promoteVendor := `
        INSERT INTO
            Vendor (vendorId, job)
        VALUES
            (
                (SELECT applicantId FROM Application WHERE id = ?),
                (SELECT job FROM Application WHERE id = ?)
            )
    `

	if _, err := tx.ExecContext(ctx, promoteVendor, id, id); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.New("Failed to promote applicant to vendor and rollback transaction")
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.New("Failed to commit transaction and rollback")
		}

		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (m *Mysql) RejectApplication(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE Application SET status = 'rejected' WHERE id = ?"

	_, err := m.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
