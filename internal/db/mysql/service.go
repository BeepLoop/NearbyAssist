package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/types"
	"time"
)

func (m *Mysql) FindServiceById(id int) (*response.ServiceDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id as serviceId,
            description,
            format(rate, 2) as rate,
            latitude, 
            longitude
        FROM 
            Service
        WHERE
            id = ?
    `

	service := &response.ServiceDetails{}
	if err := m.Conn.GetContext(ctx, service, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return service, nil
}

func (m *Mysql) FindServiceByVendor(id int) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id,
            vendorId,
            description,
            rate,
            latitude,
            longitude
        FROM 
            Service
        WHERE
            vendorId = ?
    `

	services := make([]*models.ServiceModel, 0)
	err := m.Conn.SelectContext(ctx, &services, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (m *Mysql) FindAllService() ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id,
            vendorId,
            description,
            rate,
            latitude,
            longitude
        FROM 
            Service
        LIMIT
            10
    `

	services := make([]*models.ServiceModel, 0)
	err := m.Conn.SelectContext(ctx, &services, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (m *Mysql) RegisterService(service *request.NewService) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := m.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}

	registerService := `
	        INSERT INTO
	            Service
	                (vendorId, description, rate, latitude, longitude)
	        VALUES 
                (
                    :vendorId,
                    :description,
                    :rate,
                    :latitude,
                    :longitude
                )
	    `

	res, err := tx.NamedExecContext(ctx, registerService, service)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}

		return 0, err
	}

	serviceId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	registerTag := `
        INSERT INTO 
            Service_Tag (serviceId, tagId)
        VALUES
            (
                ?,
                (SELECT id FROM Tag WHERE title = ?)
            )
    `

	var tagErr error
	for _, tag := range service.Tags {
		if _, err := tx.ExecContext(ctx, registerTag, serviceId, tag); err != nil {
			tagErr = err
			break
		}
	}

	if tagErr != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}

		return 0, err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(serviceId), nil
}

func (m *Mysql) UpdateService(service *request.UpdateService) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// TODO: Update tags and use transaction

	query := `
        UPDATE
            Service
        SET
            title = :title,
            description = :description,
            rate = :rate,
            latitude = :latitude
            longitude = :longitude
        WHERE
            id = :id
    `

	_, err := m.Conn.NamedExecContext(ctx, query, service)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (m *Mysql) DeleteService(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "DELETE FROM Service WHERE id = ?"

	_, err := m.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (m *Mysql) GeoSpatialSearch(params *types.SearchParams) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            s.id,
            s.vendorId,
            s.description,
            format(s.rate, 2) as rate,
            s.latitude,
            s.longitude
        FROM 
            Service_Tag st
            JOIN Service s ON s.id = st.serviceId
        WHERE
            st.tagId = (SELECT id from Tag WHERE title = ?)
        AND
            ST_Distance_Sphere(
                POINT(s.longitude, s.latitude),
                POINT(?, ?)
            ) < ?
    `

	services := make([]*models.ServiceModel, 0)
	err := m.Conn.SelectContext(ctx, &services, query, params.Query, params.Longitude, params.Latitude, params.Radius)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}
