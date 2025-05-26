package service_repo

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlServiceRepository struct {
	db *sqlx.DB
}

func NewMysqlServiceRepository(db *sqlx.DB) *MysqlServiceRepository {
	return &MysqlServiceRepository{
		db: db,
	}
}

func (s *MysqlServiceRepository) Create(service *models.ServiceModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	createService := `
	        INSERT INTO
	            Service (id, vendorId, title, description, price, signature)
	        VALUES 
                (:id, :vendorId, :title, :description, :price, :signature)
    `
	service.Id = utils.GenerateId()
	if _, err := tx.NamedExecContext(ctx, createService, service); err != nil {
		return "", err
	}

	serviceAddressRelation := `
        INSERT INTO
            ServiceAddress (serviceId, addressId)
        SELECT
            ?, addressId
        FROM
            UserAddress
        WHERE
            userId = ?
            
    `
	if _, err := tx.ExecContext(ctx, serviceAddressRelation, service.Id, service.VendorId); err != nil {
		return "", err
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
	for _, extra := range service.Extras {
		extraId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, registerExtra, extraId, extra.Title, extra.Description, extra.Price, service.Id); err != nil {
			return "", err
		}

		if _, err := tx.ExecContext(ctx, registerServiceExtra, service.Id, extraId); err != nil {
			return "", err
		}
	}

	tagExists := `
        SELECT EXISTS
            (SELECT 1 FROM Tag WHERE title = ?)
        AS tag_exists
    `
	getTagByTitle := "SELECT id, title FROM Tag WHERE title = ?"
	createTag := "INSERT INTO Tag (id, title) VALUES (?, ?)"
	serviceTagRelation := `
        INSERT INTO
            ServiceTag (id, serviceId, tagId)
        VALUES
            (?, ?, ?)
    `

	for _, tag := range service.TagsAsString {
		exists := false
		if err := tx.GetContext(ctx, &exists, tagExists, tag); err != nil {
			fmt.Println("error existence check: ", err.Error())
			return "", err
		}

		if exists {
			existingTag := new(models.TagModel)
			if err := tx.GetContext(ctx, existingTag, getTagByTitle, tag); err != nil {
				fmt.Println("error get existing tag: ", err.Error())
				return "", err
			}

			serviceTagId := utils.GenerateId()
			if _, err := tx.ExecContext(ctx, serviceTagRelation, serviceTagId, service.Id, existingTag.Id); err != nil {
				fmt.Println("error service tag relation create: ", err.Error())
				return "", err
			}

			continue
		}

		tagId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, createTag, tagId, tag); err != nil {
			return "", err
		}

		serviceTagId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, serviceTagRelation, serviceTagId, service.Id, tagId); err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return service.Id, nil
}

func (s *MysqlServiceRepository) CreateWithPricingType(service *models.ServiceModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	createService := `
	        INSERT INTO
	            Service (id, vendorId, title, description, price, pricingType, signature)
	        VALUES 
                (:id, :vendorId, :title, :description, :price, :pricingType, :signature)
    `
	service.Id = utils.GenerateId()
	if _, err := tx.NamedExecContext(ctx, createService, service); err != nil {
		return "", err
	}

	serviceAddressRelation := `
        INSERT INTO
            ServiceAddress (serviceId, addressId)
        SELECT
            ?, addressId
        FROM
            UserAddress
        WHERE
            userId = ?
            
    `
	if _, err := tx.ExecContext(ctx, serviceAddressRelation, service.Id, service.VendorId); err != nil {
		return "", err
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
	for _, extra := range service.Extras {
		extraId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, registerExtra, extraId, extra.Title, extra.Description, extra.Price, service.Id); err != nil {
			return "", err
		}

		if _, err := tx.ExecContext(ctx, registerServiceExtra, service.Id, extraId); err != nil {
			return "", err
		}
	}

	tagExists := `
        SELECT EXISTS
            (SELECT 1 FROM Tag WHERE title = ?)
        AS tag_exists
    `
	getTagByTitle := "SELECT id, title FROM Tag WHERE title = ?"
	createTag := "INSERT INTO Tag (id, title) VALUES (?, ?)"
	serviceTagRelation := `
        INSERT INTO
            ServiceTag (id, serviceId, tagId)
        VALUES
            (?, ?, ?)
    `

	for _, tag := range service.TagsAsString {
		exists := false
		if err := tx.GetContext(ctx, &exists, tagExists, tag); err != nil {
			fmt.Println("error existence check: ", err.Error())
			return "", err
		}

		if exists {
			existingTag := new(models.TagModel)
			if err := tx.GetContext(ctx, existingTag, getTagByTitle, tag); err != nil {
				return "", err
			}

			serviceTagId := utils.GenerateId()
			if _, err := tx.ExecContext(ctx, serviceTagRelation, serviceTagId, service.Id, existingTag.Id); err != nil {
				return "", err
			}

			continue
		}

		tagId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, createTag, tagId, tag); err != nil {
			return "", err
		}

		serviceTagId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, serviceTagRelation, serviceTagId, service.Id, tagId); err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return service.Id, nil
}

func (s *MysqlServiceRepository) FindById(serviceId string) (*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	service := new(models.ServiceModel)

	query := `
        SELECT
            id,
            vendorId,
            title,
            description,
            price,
            pricingType,
            disabled,
            status,
            rejectReason,
            createdAt,
            updatedAt,
            acceptedAt,
            rejectedAt
        FROM 
            Service
        WHERE
            id = ?
    `
	if err := s.db.GetContext(ctx, service, query, serviceId); err != nil {
		return nil, err
	}

	if address, err := s.getAddress(serviceId); err != nil {
		return nil, err
	} else {
		service.Address = *address
	}

	if extras, err := s.getExtras(serviceId); err != nil {
		return nil, err
	} else {
		service.Extras = extras
	}

	if tags, err := s.getTags(serviceId); err != nil {
		return nil, err
	} else {
		service.Tags = tags
	}

	if images, err := s.getPhotos(serviceId); err != nil {
		return nil, err
	} else {
		service.Images = images
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return service, nil
}

func (s *MysqlServiceRepository) FindBySignature(signature string) (*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id FROM Service WHERE signature = ?"
	var id string
	if err := s.db.GetContext(ctx, &id, query, signature); err != nil {
		return nil, err
	}

	service, err := s.FindById(id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return service, nil
}

func (s *MysqlServiceRepository) GetAll(limit, offset int) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id
        FROM 
            Service
        WHERE
            disabled = 0
        ORDER BY
            createdAt, updatedAt DESC
        LIMIT ? OFFSET ?
    `
	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query, limit, offset); err != nil {
		return nil, err
	}

	services := make([]*models.ServiceModel, 0)
	for _, id := range ids {
		service, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceRepository) GetAllUnderReview(limit, offset int) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id
        FROM 
            Service
        WHERE
            disabled = 0 AND status = 'under_review'
        ORDER BY
            createdAt, updatedAt ASC
        LIMIT ? OFFSET ?
    `
	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query, limit, offset); err != nil {
		return nil, err
	}

	services := make([]*models.ServiceModel, 0)
	for _, id := range ids {
		service, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceRepository) GetAllWithTag(tag string) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            s.id
        FROM 
            ServiceTag st
            JOIN Service s ON s.id = st.serviceId
            JOIN Tag t ON t.id = st.tagId
        WHERE
            t.title = ? AND s.disabled = 0 AND s.status = 'accepted'
    `
	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query, tag); err != nil {
		return nil, err
	}

	services := make([]*models.ServiceModel, 0)
	for _, id := range ids {
		service, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceRepository) GetAllWithTagAny(tags []string) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	base := `
        SELECT DISTINCT
            s.id
        FROM 
            ServiceTag st
            JOIN Service s ON s.id = st.serviceId
            JOIN Tag t ON t.id = st.tagId
        WHERE
            s.disabled = 0 AND s.status = 'accepted' AND
    `
	args := make([]interface{}, 0)
	query := base + " t.title IN ("

	placeholders := make([]string, 0)
	for _, tag := range tags {
		placeholders = append(placeholders, "?")
		args = append(args, tag)
	}

	query += strings.Join(placeholders, ", ") + ")"
	fmt.Println(query)

	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query, args...); err != nil {
		return nil, err
	}

	services := make([]*models.ServiceModel, 0)
	for _, id := range ids {
		service, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceRepository) FuzzyMatchTags(tags []string) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT DISTINCT
            s.id
        FROM 
            ServiceTag st
            JOIN Service s ON s.id = st.serviceId
            JOIN Tag t ON t.id = st.tagId
        WHERE
            s.disabled = 0 AND s.status = 'accepted'
    `
	for _, tag := range tags {
		query += fmt.Sprintf(" AND t.title LIKE '%%%s%%'", tag)
	}

	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query); err != nil {
		return nil, err
	}

	services := make([]*models.ServiceModel, 0)
	for _, id := range ids {
		service, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceRepository) GetAllTopRated(limit int) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            s.id
        FROM
            Service s
            JOIN Vendor v ON v.vendorId = s.vendorId
        ORDER BY
            v.rating DESC
        LIMIT ?
    `

	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query, limit); err != nil {
		return nil, err
	}

	services := make([]*models.ServiceModel, 0)
	for _, id := range ids {
		service, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceRepository) IsVendor(vendorId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT EXISTS (SELECT 1 FROM Vendor WHERE vendorId = ?) AS is_vendor"
	isVendor := false
	if err := s.db.GetContext(ctx, &isVendor, query, vendorId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isVendor, nil
}

func (s *MysqlServiceRepository) GetReviews(serviceId string) ([]*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            r.*
        FROM
            Review r
            JOIN Booking t ON t.id = r.bookingId
            JOIN Service s ON s.id = t.serviceId
        WHERE
            s.id = ?
        ORDER BY
            r.createdAt DESC
    `

	reviews := make([]*models.ReviewModel, 0)
	if err := s.db.SelectContext(ctx, &reviews, query, serviceId); err != nil {
		return nil, err
	}

	getRevieweeQuery := `
        SELECT
            u.id,
            u.name,
            u.email,
            u.imageUrl
        FROM
            Review r
            JOIN User u ON r.revieweeId = u.id
        WHERE
            r.id = ?
    `
	for _, review := range reviews {
		user := new(models.UserModel)
		if err := s.db.GetContext(ctx, user, getRevieweeQuery, review.Id); err != nil {
			return nil, err
		}
		review.Reviewee = user
	}

	return reviews, nil
}

func (s *MysqlServiceRepository) FindPhotoById(imageId string) (*models.ServicePhotoModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM ServicePhoto where id = ?"

	photo := new(models.ServicePhotoModel)
	if err := s.db.GetContext(ctx, photo, query, imageId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return photo, nil
}

func (s *MysqlServiceRepository) Update(updatedData *models.ServiceModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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
            price = :price,
            pricingType = :pricingType
        WHERE
            id = :id
    `
	if _, err = tx.NamedExecContext(ctx, updateService, updatedData); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	getCurrentTags := `
        SELECT
            st.id,
            st.serviceId,
            st.tagId,
            t.title
        FROM
            ServiceTag st
            JOIN Tag t ON t.id = st.tagId
        WHERE
            st.serviceId = ?
    `

	// Retrieve current tags
	currentSvcTag := make([]models.ServiceTagModel, 0)
	if err := tx.SelectContext(ctx, &currentSvcTag, getCurrentTags, updatedData.Id); err != nil {
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
	tagExists := `
        SELECT EXISTS
            (SELECT 1 FROM Tag WHERE title = ?)
        AS tag_exists
    `
	getTagByTitle := "SELECT id, title FROM Tag WHERE title = ?"
	createTag := "INSERT INTO Tag (id, title) VALUES (?, ?)"
	serviceTagRelation := `
        INSERT INTO
            ServiceTag (id, serviceId, tagId)
        VALUES
            (?, ?, ?)
    `

	for _, tag := range updatedData.TagsAsString {
		exists := false
		if err := tx.GetContext(ctx, &exists, tagExists, tag); err != nil {
			return err
		}

		if exists {
			existingTag := new(models.TagModel)
			if err := tx.GetContext(ctx, existingTag, getTagByTitle, tag); err != nil {
				return err
			}

			serviceTagId := utils.GenerateId()
			if _, err := tx.ExecContext(ctx, serviceTagRelation, serviceTagId, updatedData.Id, existingTag.Id); err != nil {
				return err
			}

			continue
		}

		tagId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, createTag, tagId, tag); err != nil {
			return err
		}

		serviceTagId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, serviceTagRelation, serviceTagId, updatedData.Id, tagId); err != nil {
			return err
		}
	}

	// Commit booking
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

func (s *MysqlServiceRepository) Resubmit(serviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Service SET status = 'under_review', updatedAt = CURRENT_TIMESTAMP() WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, serviceId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlServiceRepository) AddImage(data *models.ServicePhotoModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data.Id = utils.GenerateId()

	query := `
        INSERT INTO
            ServicePhoto (id, serviceId, vendorId, url)
        VALUES
            (:id, :serviceId, :vendorId, :url)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlServiceRepository) DeleteImage(imageId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deleteImageQuery := `
        DELETE FROM ServicePhoto
        WHERE id = ?
    `
	if _, err := s.db.ExecContext(ctx, deleteImageQuery, imageId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlServiceRepository) AddExtra(data *models.ExtraModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	data.Id = utils.GenerateId()

	insertExtraQuery := `
        INSERT INTO
            Extra (id, title, description, price, serviceId)
        VALUES 
            (:id, :title, :description, :price, :serviceId)
    `

	if _, err := tx.NamedExecContext(ctx, insertExtraQuery, data); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	insertServiceExtraQuery := `
        INSERT INTO
            ServiceExtra (serviceId, extraId)
        VALUES
            (?, ?)
    `

	if _, err := tx.ExecContext(ctx, insertServiceExtraQuery, data.ServiceId, data.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlServiceRepository) EditExtra(data *models.ExtraModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	bookingsWithThisExtraQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.cost
        FROM
            BookingExtra te
            JOIN Booking t ON t.id = te.bookingId
        WHERE
            te.extraId = ? AND (t.status = 'pending' OR t.status = 'confirmed')
    `

	bookingsWithThisExtra := make([]*models.BookingModel, 0)
	if err := tx.SelectContext(ctx, &bookingsWithThisExtra, bookingsWithThisExtraQuery, data.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	if len(bookingsWithThisExtra) != 0 {
		return errors.New("This service extra is actively in use")
	}

	updateExtraQuery := `
        UPDATE 
            Extra
        SET 
            title = :title,
            description = :description,
            price = :price
        WHERE
            id = :id
    `

	if _, err := tx.NamedExecContext(ctx, updateExtraQuery, data); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
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

func (s *MysqlServiceRepository) DeleteExtra(extraId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	markExtraAsDeletedQuery := "UPDATE Extra set deleted = 1 WHERE id = ?"
	if _, err := tx.ExecContext(ctx, markExtraAsDeletedQuery, extraId); err != nil {
		return err
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

func (s *MysqlServiceRepository) FindExtraById(extraId string) (*models.ExtraModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM Extra WHERE id = ?"

	extra := new(models.ExtraModel)
	if err := s.db.GetContext(ctx, extra, query, extraId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return extra, nil
}

func (s *MysqlServiceRepository) IsVendorRestricted(serviceId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isVendorRestrictedQuery := `
        SELECT EXISTS (
            SELECT 1 
            FROM Service s
            INNER JOIN Restricted r 
            ON s.vendorId = r.userId
            WHERE s.id = ?
        ) AS isRestricted;
    `

	isRestricted := false
	if err := s.db.GetContext(ctx, &isRestricted, isVendorRestrictedQuery, serviceId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isRestricted, nil
}

func (s *MysqlServiceRepository) IsVendorBanned(serviceId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT EXISTS (
            SELECT 1 
            FROM Service s
            INNER JOIN Ban b 
            ON s.vendorId = b.userId
            WHERE s.id = ?
        ) AS isRestricted;
    `

	isBanned := false
	if err := s.db.GetContext(ctx, &isBanned, query, serviceId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isBanned, nil
}

func (s *MysqlServiceRepository) Disable(serviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Service SET disabled = 1 WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, serviceId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlServiceRepository) Enable(serviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Service SET disabled = 0 WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, serviceId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlServiceRepository) HasActiveBookingWithThisExtra(extraId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT EXISTS
            (
                SELECT
                    1
                FROM
                    BookingExtra be
                    JOIN Booking b ON b.id = be.bookingId
                WHERE
                    be.extraId = ? AND (b.status = 'pending' OR b.status = 'confirmed')
            )
        AS has_booking
    `

	hasBooking := false
	if err := s.db.GetContext(ctx, &hasBooking, query, extraId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return hasBooking, nil
}

func (s *MysqlServiceRepository) HasActiveBookingWithThisService(serviceId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT EXISTS
            (SELECT 1 FROM Booking WHERE id = ? AND (status = 'pending' OR status = 'confirmed'))
        AS has_booking
    `
	hasBooking := false
	if err := s.db.GetContext(ctx, &hasBooking, query, serviceId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return hasBooking, nil
}

func (s *MysqlServiceRepository) Accept(serviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Service SET status = 'accepted', acceptedAt = CURRENT_TIMESTAMP() WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, serviceId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlServiceRepository) Reject(serviceId, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Service SET status = 'rejected', rejectReason = ?, rejectedAt = CURRENT_TIMESTAMP() WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, reason, serviceId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlServiceRepository) getTags(serviceId string) ([]*models.TagModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.id,
            t.title
        FROM
            ServiceTag st
            JOIN Tag t ON t.id = st.tagId
        WHERE
            st.serviceId = ?;
    `

	tags := make([]*models.TagModel, 0)
	if err := s.db.SelectContext(ctx, &tags, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return tags, nil
}

func (s *MysqlServiceRepository) getPhotos(serviceId string) ([]*models.ServicePhotoModel, error) {
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

func (s *MysqlServiceRepository) getAddress(serviceId string) (*models.AddressModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            a.id, a.address, a.latitude, a.longitude
        FROM
            Address a
            JOIN ServiceAddress sa ON sa.addressId = a.id
        WHERE
            sa.serviceId = ?
    `
	address := new(models.AddressModel)
	if err := s.db.GetContext(ctx, address, query, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return address, nil
}

func (s *MysqlServiceRepository) getExtras(serviceId string) ([]*models.ExtraModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
            se.serviceId = ? AND e.deleted = 0
    `

	extras := make([]*models.ExtraModel, 0)
	if err := s.db.SelectContext(ctx, &extras, extrasQuery, serviceId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return extras, nil
}
