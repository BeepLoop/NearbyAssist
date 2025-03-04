package vendor_repo

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlVendorRepository struct {
	db *sqlx.DB
}

func NewMysqlVendorRepository(db *sqlx.DB) *MysqlVendorRepository {
	return &MysqlVendorRepository{db: db}
}

func (s *MysqlVendorRepository) FindById(id string) (*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	vendor := new(models.VendorModel)
	query := `
        SELECT  
            v.vendorId AS id,
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
	if err := s.db.GetContext(ctx, vendor, query, id); err != nil {
		return nil, err
	}

	if restricted, err := s.IsRestricted(vendor.Id); err != nil {
		return nil, err
	} else {
		vendor.Restricted = restricted
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
	if err := s.db.SelectContext(ctx, &expertise, expertiseQuery, id); err != nil {
		return nil, err
	}
	vendor.Expertise = expertise

	getSocialsQuery := `
        SELECT
            url
        FROM
            Social
        WHERE
            userId = ?
    `

	socials := make([]string, 0)
	if err := s.db.SelectContext(ctx, &socials, getSocialsQuery, id); err != nil {
		return nil, err
	}
	vendor.Socials = socials

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}

func (s *MysqlVendorRepository) GetVendorServiceList(vendorId string) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id,
            vendorId,
            title,
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

func (s *MysqlVendorRepository) GetTags(serviceId string) ([]*models.TagModel, error) {
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

func (s *MysqlVendorRepository) IsRestricted(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            CASE
                WHEN EXISTS (SELECT 1 FROM Restricted WHERE userId = ?)
                THEN 1
                ELSE 0
            END AS user_exists;
    `

	isRestricted := false
	if err := s.db.GetContext(ctx, &isRestricted, query, userId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isRestricted, nil
}
