package models

import (
	"context"
	"fmt"
	"nearbyassist/internal/id_generator"
	"nearbyassist/internal/utils"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type ServiceSearchResult struct {
	ServiceModel
	Vendor string `json:"vendor" db:"vendor"`
}

type ServiceModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	VendorId    string `json:"vendorId" db:"vendorId" validate:"required"`
	Description string `json:"description" db:"description" validate:"required"`
	Rate        string `json:"rate" db:"rate" validate:"required"`

	// Additional fields for joins
	Tags []string `json:"tags" db:"tags" validate:"required"`
}

func NewServiceModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *ServiceModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &ServiceModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func NewServiceModelWithId(id string, conn *sqlx.DB) *ServiceModel {
	return &ServiceModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func (s *ServiceModel) EncryptDescription(encryptFunc func(string) (string, error)) (*ServiceModel, error) {
	if encrypted, err := encryptFunc(s.Description); err != nil {
		return nil, err
	} else {
		s.Description = encrypted
	}
	return s, nil
}

func (s *ServiceModel) DecryptDescription(decryptFunc func(string) (string, error)) (*ServiceModel, error) {
	if decrypted, err := decryptFunc(s.Description); err != nil {
		return nil, err
	} else {
		s.Description = decrypted
	}
	return s, nil
}

func (s *ServiceModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	registerService := `
	        INSERT INTO
	            Service
	                (id, vendorId, description, rate, latitude, longitude)
	        VALUES 
                (
                    :id,
                    :vendorId,
                    :description,
                    :rate,
                    :latitude,
                    :longitude
                )
	    `

	if _, err := tx.NamedExecContext(ctx, registerService, s); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	registerTag := `
        INSERT INTO 
            ServiceTag (serviceId, tagId)
        VALUES
            (
                ?,
                (SELECT id FROM Tag WHERE title = ?)
            )
    `

	var tagErr error
	for _, tag := range s.Tags {
		if _, err := tx.ExecContext(ctx, registerTag, s.Id, tag); err != nil {
			tagErr = err
			break
		}
	}

	if tagErr != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return s.Id, nil
}

func (s *ServiceModel) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM Service"

	count := 0
	err := s.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *ServiceModel) FindById(id string) (*ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id,
            description,
            format(rate, 2) as rate,
            latitude, 
            longitude
        FROM 
            Service
        WHERE
            id = ?
    `

	if err := s.Conn.GetContext(ctx, s, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return s, nil
}

func (s *ServiceModel) FindAll() ([]*ServiceModel, error) {
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

	services := make([]*ServiceModel, 0)
	err := s.Conn.SelectContext(ctx, &services, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *ServiceModel) FindByAllVendorId(id string) ([]*ServiceModel, error) {
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

	services := make([]*ServiceModel, 0)
	if err := s.Conn.SelectContext(ctx, &services, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *ServiceModel) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.Conn.BeginTxx(ctx, nil)
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

	if _, err = tx.NamedExecContext(ctx, updateService, s); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	currentServiceTagsQuery := `
        SELECT
            st.id,
            st.serviceId,
            t.title AS tag
        FROM
            ServiceTag st
            JOIN Tag t ON t.id = st.tagId
        WHERE
            st.serviceId = ?;
    `

	currentServiceTags := make([]ServiceTagModel, 0)
	if err := tx.SelectContext(ctx, &currentServiceTags, currentServiceTagsQuery, s.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	deleteTag := "DELETE FROM ServiceTag WHERE id = ?"
	insertTag := "INSERT INTO ServiceTag (serviceId, tagId) VALUES (?, (SELECT id FROM Tag WHERE title = ?))"

	newTags := s.Tags

	for _, tag := range currentServiceTags {
		exists := utils.StringSliceContains(newTags, tag.TagId)
		if exists {
			// Remove item from newTags
			newTags = utils.RemoveStringFromSlice(newTags, tag.TagId)
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
		if _, err := tx.ExecContext(ctx, insertTag, s.Id, tag); err != nil {
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

func (s *ServiceModel) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "DELETE FROM Service WHERE id = ?"
	if _, err := s.Conn.ExecContext(ctx, query, s.Id); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *ServiceModel) GeoSpatialSearch(params map[string]string) ([]*ServiceSearchResult, error) {
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
            ServiceTag st
            JOIN Service s ON s.id = st.serviceId
            JOIN User u ON u.id = s.vendorId
        WHERE
    `

	if q, ok := params["q"]; ok {
		condition := ""
		tags := strings.Split(q, ",")
		for i, tag := range tags {
			if i == 0 {
				condition += fmt.Sprintf(" st.tagId = (SELECT id from Tag WHERE title = '%s')", tag)
			} else {
				condition += fmt.Sprintf(" OR st.tagId = (SELECT id from Tag WHERE title = '%s')", tag)
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
		condition := fmt.Sprintf(" AND ST_Distance_Sphere(POINT(s.longitude, s.latitude), POINT(%v, %v))", location[0], location[1])
		query += condition
	} else {
		return nil, fmt.Errorf("Missing location parameter 'l'")
	}

	if r, ok := params["r"]; ok {
		query += fmt.Sprintf(" < %v", r)
	} else {
		return nil, fmt.Errorf("Missing radius parameter 'r'")
	}

	services := make([]*ServiceSearchResult, 0)
	err := s.Conn.SelectContext(ctx, &services, query)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *ServiceModel) GetTags() ([]string, error) {
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
	if err := s.Conn.SelectContext(ctx, &tags, query, s.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return tags, nil
}

func (s *ServiceModel) GetPhotos() ([]ServicePhotoModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, url FROM ServicePhoto WHERE serviceId = ?"

	images := make([]ServicePhotoModel, 0)
	if err := s.Conn.SelectContext(ctx, &images, query, s.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return images, nil
}

func (s *ServiceModel) GetReviews() ([]ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, serviceId, rating FROM Review WHERE serviceId = ?"

	reviews := make([]ReviewModel, 0)
	if err := s.Conn.SelectContext(ctx, &reviews, query, s.Id); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *ServiceModel) GetVendor() (*VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            s.vendorId,
            u.name as vendor
        FROM 
            Service s
        JOIN 
            User u ON u.id = s.vendorId
        WHERE
            s.id = ?
    `

	vendor := new(VendorModel)
	if err := s.Conn.GetContext(ctx, vendor, query, s.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}

//
// func StringSliceContains(slice []string, target string) bool {
// 	for _, item := range slice {
// 		if item == target {
// 			return true
// 		}
// 	}
//
// 	return false
// }
//
// func RemoveStringFromSlice(slice []string, target string) []string {
// 	var result []string
//
// 	for _, item := range slice {
// 		if item != target {
// 			result = append(result, item)
// 		}
// 	}
//
// 	return result
// }
