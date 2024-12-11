package service_repo

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type MysqlServiceRepository struct {
	db *sqlx.DB
}

func NewMysqlServiceRepository(db *sqlx.DB) *MysqlServiceRepository {
	return &MysqlServiceRepository{
		db: db,
	}
}

func (s *MysqlServiceRepository) Create(data *models.ServiceModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	if id, err := gonanoid.New(); err != nil {
		return "", errors.New("Failed to generate id for service")
	} else {
		data.Id = id
	}

	registerService := `
	        INSERT INTO
	            Service
	                (id, vendorId, title, description, rate, latitude, longitude, signature)
	        VALUES 
                (
                    :id,
                    :vendorId,
                    :title,
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
		tagId, err := gonanoid.New()
		if err != nil {
			return "", errors.New("Failed to generate id for tag")
		}

		if _, err := tx.ExecContext(ctx, registerTag, tagId, data.Id, tag); err != nil {
			return "", err
		}
	}

	registerExtra := `
        INSERT INTO
            Extra (id, title, description, price, serviceId)
        VALUES 
            (?, ?, ?, ?, ?)
    `
	registerServiceExtra := `
        INSERT INTO 
            ServiceExtra (serviceId, extraId)
        VALUES
            (?, ?)
    `
	for _, extra := range data.Extras {
		extraId, err := gonanoid.New()
		if err != nil {
			return "", errors.New("Failed to generate id for service extra")
		}

		if _, err := tx.ExecContext(ctx, registerExtra, extraId, extra.Title, extra.Description, extra.Price, data.Id); err != nil {
			return "", err
		}

		if _, err := tx.ExecContext(ctx, registerServiceExtra, data.Id, extraId); err != nil {
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

func (s *MysqlServiceRepository) FindAll() ([]*models.ServiceModel, error) {
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

func (s *MysqlServiceRepository) FindById(id string) (*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	service := new(models.ServiceModel)

	query := `
        SELECT
            id,
            vendorId,
            title,
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

	extrasQuery := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN ServiceExtra se ON se.extraId = e.id
        WHERE
            se.serviceId = ?
    `

	extras := make([]models.ExtraModel, 0)
	if err := s.db.SelectContext(ctx, &extras, extrasQuery, id); err != nil {
		return nil, err
	}
	service.Extras = extras

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return service, nil
}

func (s *MysqlServiceRepository) FindBySignature(signature string) (*models.ServiceModel, error) {
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

func (s *MysqlServiceRepository) IsVendor(vendorId string) error {
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

func (s *MysqlServiceRepository) GetVendorInfo(vendorId string) (*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	vendor := new(models.VendorModel)
	query := `
        SELECT  
            v.vendorId,
            v.rating,
            v.restricted,
            u.name AS vendor,
            u.email AS email,
            u.imageUrl AS imageUrl
        FROM 
            Vendor  v
            JOIN User u ON u.id = v.vendorId
        WHERE 
            vendorId = ?
    `
	if err := s.db.GetContext(ctx, vendor, query, vendorId); err != nil {
		return nil, err
	}

	expertiseQuery := `
        SELECT
            e.title
        FROM
            Expertise e
            JOIN VendorExpertise ve ON ve.expertiseId = e.id
        WHERE
            ve.vendorId = ?
    `

	expertise := make([]string, 0)
	if err := s.db.SelectContext(ctx, &expertise, expertiseQuery, vendorId); err != nil {
		return nil, err
	}

	vendor.Expertise = expertise

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}

func (s *MysqlServiceRepository) GetTags(serviceId string) ([]string, error) {
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

func (s *MysqlServiceRepository) GetReviews(serviceId string) ([]*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, serviceId, rating FROM Review WHERE serviceId = ?"

	reviews := make([]*models.ReviewModel, 0)
	if err := s.db.SelectContext(ctx, &reviews, query, serviceId); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *MysqlServiceRepository) GetPhotos(serviceId string) ([]*models.ServicePhotoModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, serviceId, url, vendorId FROM ServicePhoto WHERE serviceId = ?"

	images := make([]*models.ServicePhotoModel, 0)
	if err := s.db.SelectContext(ctx, &images, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return images, nil
}

func (s *MysqlServiceRepository) Update(updatedService *models.ServiceModel) error {
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
            title = :title,
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
		generatedId, err := gonanoid.New()
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

func (s *MysqlServiceRepository) Delete(serviceId string) error {
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

func (s *MysqlServiceRepository) GetAllByVendorId(vendorId string) ([]*models.ServiceModel, error) {
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

func (s *MysqlServiceRepository) GeoSpatialSearch(params map[string]string) ([]*models.GeoSpatialSearchResult, error) {
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
