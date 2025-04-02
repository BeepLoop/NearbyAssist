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

func (s *MysqlServiceRepository) Create(data *models.ServiceModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	data.Id = utils.GenerateId()

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
	for _, tag := range data.TagsAsString {
		tagId := utils.GenerateId()
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
		extraId := utils.GenerateId()
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

func (s *MysqlServiceRepository) FindAll(limit, offset int) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id,
            vendorId,
            title,
            description,
            FORMAT(rate, 2) AS rate,
            latitude,
            longitude,
            createdAt
        FROM 
            Service
        ORDER BY
            createdAt, updatedAt DESC
        LIMIT ?
        OFFSET ?
    `

	services := make([]*models.ServiceModel, 0)
	if err := s.db.SelectContext(ctx, &services, query, limit, offset); err != nil {
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
            se.serviceId = ? AND e.deleted = 0
    `

	for _, service := range services {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, extrasQuery, service.Id); err != nil {
			return nil, err
		}
		service.Extras = extras

		if tags, err := s.GetTags(service.Id); err != nil {
			return nil, err
		} else {
			service.Tags = tags
		}

		if images, err := s.GetPhotos(service.Id); err != nil {
			return nil, err
		} else {
			service.Images = images
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlServiceRepository) FindAllByTag(tag string) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	services := make([]*models.ServiceModel, 0)

	query := `
        SELECT
            s.id,
            s.createdAt,
            s.vendorId,
            s.title,
            s.description,
            FORMAT(s.rate, 2) as rate,
            s.latitude, 
            s.longitude
        FROM 
            ServiceTag st
            JOIN Service s ON s.id = st.serviceId
            JOIN Tag t ON t.id = st.tagId
        WHERE
            t.title = ?
    `

	if err := s.db.SelectContext(ctx, &services, query, tag); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
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
            format(rate, 2) as rate,
            latitude, 
            longitude
        FROM 
            Service
        WHERE
            id = ?
    `
	if err := s.db.GetContext(ctx, service, query, serviceId); err != nil {
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
            se.serviceId = ? AND e.deleted = 0
    `

	extras := make([]*models.ExtraModel, 0)
	if err := s.db.SelectContext(ctx, &extras, extrasQuery, serviceId); err != nil {
		return nil, err
	}
	service.Extras = extras

	if tags, err := s.GetTags(serviceId); err != nil {
		return nil, err
	} else {
		service.Tags = tags
	}

	if images, err := s.GetPhotos(serviceId); err != nil {
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

	isVendor := false
	query := "SELECT EXISTS (SELECT 1 FROM Vendor WHERE vendorId = ?) AS is_vendor"
	if err := s.db.GetContext(ctx, &isVendor, query, vendorId); err != nil {
		return err
	}

	if !isVendor {
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
            u.name AS vendor,
            u.email AS email,
            u.phone AS phone,
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

	if vendor.VendorId == "" {
		return nil, errors.New("not found")
	}

	if restricted, err := s.IsVendorRestricted(vendor.VendorId); err != nil {
		return nil, err
	} else {
		vendor.Restricted = restricted
	}

	expertiseQuery := `
        SELECT
            e.title
        FROM
            Expertise e
            JOIN UserExpertise ve ON ve.expertiseId = e.id
        WHERE
            ve.userId = ?
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

func (s *MysqlServiceRepository) GetTags(serviceId string) ([]*models.TagModel, error) {
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

func (s *MysqlServiceRepository) GetReviews(serviceId string) ([]*models.ReviewModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// TODO: Implement this, join on transaction and review table

	query := ""

	reviews := make([]*models.ReviewModel, 0)
	if err := s.db.SelectContext(ctx, &reviews, query, serviceId); err != nil {
		return nil, err
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

	checkIfHasActiveTransactionsQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.cost
        FROM
            Transaction t
        WHERE
            t.id = ? AND (t.status = 'pending' OR t.status = 'confirmed')
    `
	activeTransactions := make([]*models.TransactionModel, 0)
	if err := tx.SelectContext(ctx, &activeTransactions, checkIfHasActiveTransactionsQuery, updatedService.Id); err != nil {
		return err
	}

	if len(activeTransactions) != 0 {
		return errors.New("This service is actively in use")
	}

	updateService := `
        UPDATE
            Service
        SET
            title = :title,
            description = :description,
            rate = :rate
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
	for _, tag := range updatedService.TagsAsString {
		generatedId := utils.GenerateId()
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

	transactionsWithThisExtraQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.cost
        FROM
            TransactionExtra te
            JOIN Transaction t ON t.id = te.transactionId
        WHERE
            te.extraId = ? AND (t.status = 'pending' OR t.status = 'confirmed')
    `

	transactionsWithThisExtra := make([]*models.TransactionModel, 0)
	if err := tx.SelectContext(ctx, &transactionsWithThisExtra, transactionsWithThisExtraQuery, data.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	if len(transactionsWithThisExtra) != 0 {
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

	transactionsWithThisExtraQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.cost
        FROM
            TransactionExtra te
            JOIN Transaction t ON t.id = te.transactionId
        WHERE
            te.extraId = ? AND (t.status = 'pending' OR t.status = 'confirmed')
    `

	transactionsWithThisExtra := make([]*models.TransactionModel, 0)
	if err := tx.SelectContext(ctx, &transactionsWithThisExtra, transactionsWithThisExtraQuery, extraId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	if len(transactionsWithThisExtra) != 0 {
		return errors.New("This service extra is actively in use")
	}

	markExtraAsDeletedQuery := "UPDATE Extra set deleted = 1 WHERE id = ?"
	if _, err := tx.ExecContext(ctx, markExtraAsDeletedQuery, extraId); err != nil {
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
				condition += fmt.Sprintf(" t.title LIKE '%%%s%%'", cleaned)
			} else {
				condition += fmt.Sprintf(" OR t.title LIKE '%%%s%%'", cleaned)
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
