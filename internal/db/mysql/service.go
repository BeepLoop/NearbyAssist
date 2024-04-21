package mysql

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/types"
	"time"
)

func (m *Mysql) FindServiceById(id int) (*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id, vendorId, title, description, rate, ST_AsText(location) as location, categoryId
        FROM 
            Service
        WHERE
            id = ?
    `

	service := models.NewServiceModel()
	err := m.Conn.GetContext(ctx, service, query, id)
	if err != nil {
		return nil, err
	}

	err = models.ExtractLatLongFromLocation(&service.GeoSpatialModel)
	if err != nil {
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
            id, vendorId, title, description, rate, ST_AsText(location) as location, categoryId
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

	for _, service := range services {
		err := models.ExtractLatLongFromLocation(&service.GeoSpatialModel)
		if err != nil {
			return nil, err
		}
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
            id, vendorId, title, description, rate, ST_AsText(location) as location, categoryId
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

	for _, service := range services {
		err := models.ExtractLatLongFromLocation(&service.GeoSpatialModel)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (m *Mysql) RegisterService(service *models.ServiceModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	        INSERT INTO
	            Service
	                (vendorId, title, description, rate, location, categoryId)
	        VALUES 
                (
                    :vendorId,
                    :title,
                    :description,
                    :rate,
                    ST_GeomFromText(:location, 4326),
                    :categoryId
                )
	    `

	models.ConstructLocationFromLatLong(&service.GeoSpatialModel)

	res, err := m.Conn.NamedExecContext(ctx, query, service)
	if err != nil {
		return -1, err
	}

	insertId, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return int(insertId), nil
}

func (m *Mysql) UpdateService(service *models.ServiceModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        UPDATE
            Service
        SET
            title = :title,
            description = :description,
            rate = :rate,
            location = ST_GeomFromText(:location, 4326),
            categoryId = :categoryId
        WHERE
            id = :id
    `

	models.ConstructLocationFromLatLong(&service.GeoSpatialModel)

	_, err := m.Conn.NamedExecContext(ctx, query, service)
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

	query := fmt.Sprintf(`
        SELECT
            id, vendorId, title, description, format(rate, 2) as rate, ST_AsText(location) as location, categoryId
        FROM 
            Service
        WHERE
            title LIKE '%%%s%%'
        AND
            ST_Distance_Sphere(
                location,
                ST_GeomFromText('POINT(%f %f)', 4326)
            ) < ?;
    `, params.Query, params.Latitude, params.Longitude)

	services := make([]*models.ServiceModel, 0)
	err := m.Conn.SelectContext(ctx, &services, query, params.Radius)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		err := models.ExtractLatLongFromLocation(&service.GeoSpatialModel)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}
