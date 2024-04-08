package models

import (
	"context"
	"errors"
	"nearbyassist/internal/db"
	"time"
)

type ApplicationModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	ApplicantId int    `json:"applicantId" db:"applicantId" validate:"required"`
	Job         string `json:"job" db:"job" validate:"required"`
	Status      string `json:"status" db:"status"`
}

func NewApplicationModel() *ApplicationModel {
	return &ApplicationModel{}
}

func (a *ApplicationModel) Create() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO 
            Application (applicantId, job, latitude, longitude)
        VALUES
            (:applicantId, :job, :latitude, :longitude)
    `

	res, err := db.Connection.NamedExecContext(ctx, query, a)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}

func (a *ApplicationModel) Update(id int) error {
	return nil
}

func (a *ApplicationModel) Delete(id int) error {
	return nil
}

func (a *ApplicationModel) FindAll(filter string) ([]*ApplicationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            id, applicantId, job, status, latitude, longitude
        FROM
            Application
    `

	switch filter {
	case "pending":
		query += " WHERE status = 'pending'"
	case "approved":
		query += " WHERE status = 'approved'"
	case "rejected":
		query += " WHERE status = 'rejected'"
	}

	applications := make([]*ApplicationModel, 0)
	err := db.Connection.SelectContext(ctx, &applications, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return applications, nil
}

func (a *ApplicationModel) FindById(applicationId int) (*ApplicationModel, error) {
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

	application := new(ApplicationModel)
	err := db.Connection.GetContext(ctx, application, query, applicationId)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return application, nil
}

func (a *ApplicationModel) Approve(applicationId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := db.Connection.Beginx()
	if err != nil {
		return err
	}

	updateStatus := `
        UPDATE 
            Application 
        SET
            status = 'approved'
        WHERE
            id = ?
    `

	if _, err := tx.ExecContext(ctx, updateStatus, applicationId); err != nil {
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

	if _, err := tx.ExecContext(ctx, promoteVendor, applicationId, applicationId); err != nil {
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

func (a *ApplicationModel) Reject(applicationId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        UPDATE
            Application 
        SET 
            status = 'rejected'
        WHERE 
            id = ?
    `

	_, err := db.Connection.ExecContext(ctx, query, applicationId)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (a *ApplicationModel) Count(filter string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            COUNT(*)
        FROM
            Application
    `

	switch filter {
	case "pending":
		query += " WHERE status = 'pending'"
	case "approved":
		query += " WHERE status = 'approved'"
	case "rejected":
		query += " WHERE status = 'rejected'"
	}

	count := 0
	err := db.Connection.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}
