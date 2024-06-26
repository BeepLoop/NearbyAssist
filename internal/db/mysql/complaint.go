package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"time"
)

func (m *Mysql) CountSystemComplaint() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM SystemComplaint"

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

func (m *Mysql) FileVendorComplaint(complaint *request.NewComplaint) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO 
            Complaint (vendorId, code, title, content)
        VALUES
            (:vendorId, :code, :title, :content)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, complaint)
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

func (m *Mysql) FileSystemComplaint(complaint *request.SystemComplaint) (int, error) {
	// TODO: Implement this function
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO 
            SystemComplaint (title, detail)
        VALUES
            (:title, :detail)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, complaint)
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

func (m *Mysql) NewSystemComplaintImage(model *models.SystemComplaintImageModel) (int, error) {
	// TODO: implement this function
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO 
            SystemComplaintImage (complaintId, url)
        VALUES
            (:complaintId, :url)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, model)
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
