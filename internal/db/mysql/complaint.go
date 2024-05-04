package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"time"
)

func (m *Mysql) CountComplaint() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM Complaint"

	count := -1
	err := m.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return count, nil
}

func (m *Mysql) FileComplaint(complaint *request.NewComplaint) (int, error) {
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
