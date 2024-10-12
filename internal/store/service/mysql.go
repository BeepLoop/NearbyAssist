package service

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlServiceStore struct {
	db *sqlx.DB
}

func NewMysqlServiceStore(db *sqlx.DB) *MysqlServiceStore {
	return &MysqlServiceStore{
		db: db,
	}
}

func (s *MysqlServiceStore) Create(data *models.ServiceModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	if id, err := store.GenerateNanoId(); err != nil {
		return "", errors.New("Failed to generate id for service")
	} else {
		data.Id = id
	}

	registerService := `
	        INSERT INTO
	            Service
	                (id, vendorId, description, rate, latitude, longitude, signature)
	        VALUES 
                (
                    :id,
                    :vendorId,
                    :description,
                    :rate,
                    :latitude,
                    :longitude,
                    :signature
                )
	    `
	if _, err := tx.NamedExecContext(ctx, registerService, data); err != nil {
		return "", err
	}

	registerTag := `
        INSERT INTO 
            ServiceTag (id, serviceId, tagId)
        VALUES
            (
                ?,
                ?,
                (SELECT id FROM Tag WHERE title = ?)
            )
    `
	for _, tag := range data.Tags {
		tagId, err := store.GenerateNanoId()
		if err != nil {
			return "", errors.New("Failed to generate id for tag")
		}

		if _, err := tx.ExecContext(ctx, registerTag, tagId, data.Id, tag); err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlServiceStore) FindAll() ([]*models.ServiceModel, error) {
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
	err := s.db.SelectContext(ctx, &services, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceStore) FindById(id string) (*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	service := new(models.ServiceModel)

	query := `
        SELECT
            id,
            vendorId,
            description,
            format(rate, 2) as rate,
            latitude, 
            longitude
        FROM 
            Service
        WHERE
            id = ?
    `
	if err := s.db.GetContext(ctx, service, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return service, nil
}

func (s *MysqlServiceStore) FindBySignature(signature string) (*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	service := new(models.ServiceModel)
	query := "SELECT id, vendorId, description, rate, latitude, longitude FROM Service WHERE signature = ?"
	if err := s.db.GetContext(ctx, service, query, signature); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return service, nil
}

func (s *MysqlServiceStore) IsVendor(vendorId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count := 0
	query := "SELECT COUNT(id) FROM Vendor WHERE vendorId = ?"
	if err := s.db.GetContext(ctx, &count, query, vendorId); err != nil {
		return err
	}

	if count == 0 {
		return errors.New("vendorId not vendor")
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlServiceStore) GetVendorInfo(vendorId string) (*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	vendor := new(models.VendorModel)

	query := `
        SELECT
            v.vendorId,
            v.rating,
            v.job,
            u.name as vendor,
            u.imageUrl as imageUrl
        FROM
            Vendor v
            JOIN User u ON u.id = v.vendorId
        WHERE 
            v.vendorId = ?
    `
	if err := s.db.GetContext(ctx, vendor, query, vendorId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}

func (s *MysqlServiceStore) GetTags(serviceId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.title AS tag
        FROM
            ServiceTag st
            JOIN Tag t ON t.id = st.tagId
        WHERE
            st.serviceId = ?;
    `

	tags := make([]string, 0)
	if err := s.db.SelectContext(ctx, &tags, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return tags, nil
}

func (s *MysqlServiceStore) GetReviews(serviceId string) ([]*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, serviceId, rating FROM Review WHERE serviceId = ?"

	reviews := make([]*models.ReviewModel, 0)
	if err := s.db.SelectContext(ctx, &reviews, query, serviceId); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *MysqlServiceStore) GetPhotos(serviceId string) ([]*models.ServicePhotoModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, url FROM ServicePhoto WHERE serviceId = ?"

	images := make([]*models.ServicePhotoModel, 0)
	if err := s.db.SelectContext(ctx, &images, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return images, nil
}

func (s *MysqlServiceStore) Update(updatedService *models.ServiceModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Start transaction
	tx, err := s.db.BeginTxx(ctx, nil)
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
	if _, err = tx.NamedExecContext(ctx, updateService, updatedService); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	getCurrentTags := `
        SELECT
            st.id,
            st.serviceId,
            t.title
        FROM
            ServiceTag st
            JOIN Tag t ON t.id = st.tagId
        WHERE
            st.serviceId = ?
    `

	// Retrieve current tags
	currentSvcTag := make([]models.ServiceTagModel, 0)
	if err := tx.SelectContext(ctx, &currentSvcTag, getCurrentTags, updatedService.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	// Delete old tags
	deleteTag := "DELETE FROM ServiceTag WHERE id = ?"
	for _, tag := range currentSvcTag {
		if _, err := tx.ExecContext(ctx, deleteTag, tag.Id); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	// Insert new tags
	insertTag := `
        INSERT INTO ServiceTag (id, serviceId, tagId) 
        SELECT ?, ?, t.id
        FROM Tag t 
        WHERE t.title = ?
    `
	for _, tag := range updatedService.Tags {
		generatedId, err := store.GenerateNanoId()
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, insertTag, generatedId, updatedService.Id, tag); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	// Commit transaction
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

func (s *MysqlServiceStore) Delete(serviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "DELETE FROM Service WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, serviceId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlServiceStore) GetAllByVendorId(vendorId string) ([]*models.ServiceModel, error) {
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
	if err := s.db.SelectContext(ctx, &services, query, vendorId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceStore) GeoSpatialSearch(params map[string]string) ([]*models.GeoSpatialSearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            s.id,
            s.vendorId,
            u.name AS vendorName,
            format(s.rate, 2) AS rate,
            s.latitude,
            s.longitude,
            (
                SELECT v.rating
                FROM Vendor v
                WHERE v.vendorId = s.vendorId
            ) AS rating,
            (
                SELECT COUNT(id)
                FROM Transaction t
                WHERE t.vendorId = s.vendorId AND t.status = 'done' AND t.serviceId = s.id
            ) AS transactions
        FROM 
            ServiceTag st
            JOIN Service s ON s.id = st.serviceId
            JOIN User u ON u.id = s.vendorId
            JOIN Tag t ON t.id = st.tagId
        WHERE
    `

	if q, ok := params["q"]; ok {
		condition := ""
		tags := strings.Split(q, ",")
		for i, tag := range tags {
			cleaned := strings.ReplaceAll(tag, "_", " ")
			if i == 0 {
				condition += fmt.Sprintf(" t.title = '%s'", cleaned)
			} else {
				condition += fmt.Sprintf(" OR t.title = '%s'", cleaned)
			}
		}

		query += condition
	} else {
		return nil, fmt.Errorf("Missing query parameter 'q'")
	}

	if l, ok := params["l"]; ok {
		location := strings.Split(l, ",")
		if len(location) != 2 {
			return nil, fmt.Errorf("Malformed location parameter 'l'")
		}
		condition := fmt.Sprintf(" AND ST_Distance_Sphere(POINT(s.longitude, s.latitude), POINT(%v, %v))", location[1], location[0])
		query += condition
	} else {
		return nil, fmt.Errorf("Missing location parameter 'l'")
	}

	if r, ok := params["r"]; ok {
		query += fmt.Sprintf(" < %v", r)
	}

	services := make([]*models.GeoSpatialSearchResult, 0)
	err := s.db.SelectContext(ctx, &services, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}
