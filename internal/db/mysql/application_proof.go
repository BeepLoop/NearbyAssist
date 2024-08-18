package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) NewApplicationProof(data *models.ApplicationProofModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO
            ApplicationProof (id, applicationId, applicantId, url)
        VALUES
            (:id, :applicationId, :applicantId, :url)
    `

	if _, err := m.Conn.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}
