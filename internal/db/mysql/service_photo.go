package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	"time"
)

func (m *Mysql) NewServicePhoto(data *models.ServicePhotoModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO
            ServicePhoto (vendorId, serviceId, url)
        VALUES
            (:vendorId, :serviceId, :url)
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

func (m *Mysql) FindAllPhotosByServiceId(serviceId int) ([]response.ServiceImages, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT 
            id AS imageId,
            url AS imageUrl
        FROM
            ServicePhoto
        WHERE
            serviceId = ?
    `

	images := make([]response.ServiceImages, 0)
	err := m.Conn.SelectContext(ctx, &images, query, serviceId)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return images, nil
}
