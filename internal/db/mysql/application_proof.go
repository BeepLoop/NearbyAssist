package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) NewApplicationProof(data *models.ApplicationProofModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO
            ApplicationProof (applicationId, applicantId, url)
        VALUES
            (:applicationId, :applicantId, :url)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, data)
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
