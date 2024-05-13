package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
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

	tx, err := m.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	updateService := `
        UPDATE
            Service
        SET
            description = :description,
            rate = :rate,
            latitude = :latitude,
            longitude = :longitude
        WHERE
            id = :id
    `

	_, err = tx.NamedExecContext(ctx, updateService, service)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	getExistingTags := `
        SELECT
            st.id,
            st.serviceId,
            t.title AS tag
        FROM
            Service_Tag st
            JOIN Tag t ON t.id = st.tagId
        WHERE
            st.serviceId = ?;
    `

	existingTags := make([]models.ServiceTagModel, 0)
	if err := tx.SelectContext(ctx, &existingTags, getExistingTags, service.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	deleteTag := "DELETE FROM Service_Tag WHERE id = ?"
	insertTag := "INSERT INTO Service_Tag (serviceId, tagId) VALUES (?, (SELECT id FROM Tag WHERE title = ?))"

	newTags := service.Tags

	for _, tag := range existingTags {
		exists := utils.StringSliceContains(newTags, tag.Tag)
		if exists {
			// Remove item from newTags
			newTags = utils.RemoveStringFromSlice(newTags, tag.Tag)
		} else {
			// Append to tagsToBeDeleted
			if _, err := tx.ExecContext(ctx, deleteTag, tag.Id); err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}

				return err
			}
		}
	}

	for _, tag := range newTags {
		if _, err := tx.ExecContext(ctx, insertTag, service.Id, tag); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

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

func (m *Mysql) FindServiceOwner(id int) (*response.ServiceOwner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            s.vendorId AS id,
            u.name as name
        FROM 
            Service s
        JOIN 
            User u ON u.id = s.vendorId
        WHERE
            s.id = ?
    `

	owner := &response.ServiceOwner{}
	if err := m.Conn.GetContext(ctx, owner, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return owner, nil
}

func (m *Mysql) GeoSpatialSearch(params *types.SearchParams) ([]*models.ServiceSearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            s.id,
            s.vendorId,
            u.name as vendor,
            s.description,
            format(s.rate, 2) as rate,
            s.latitude,
            s.longitude
        FROM 
            Service_Tag st
            JOIN Service s ON s.id = st.serviceId
            JOIN User u ON u.id = s.vendorId
        WHERE
            st.tagId = (SELECT id from Tag WHERE title = ?)
        AND
            ST_Distance_Sphere(
                POINT(s.longitude, s.latitude),
                POINT(?, ?)
            ) < ?
    `

	services := make([]*models.ServiceSearchResult, 0)
	err := m.Conn.SelectContext(ctx, &services, query, params.Query, params.Longitude, params.Latitude, params.Radius)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}
