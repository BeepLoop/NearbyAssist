package vendor_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlVendorRepository struct {
	db *sqlx.DB
}

func NewMysqlVendorRepository(db *sqlx.DB) *MysqlVendorRepository {
	return &MysqlVendorRepository{db: db}
}

func (s *MysqlVendorRepository) GetAll(limit, offset int) ([]*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT  
            v.joinedAt,
            v.vendorId,
            v.rating,
            u.name AS name,
            u.email AS email,
            u.phone AS phone,
            u.imageUrl AS imageUrl
        FROM 
            Vendor  v
            JOIN User u ON u.id = v.vendorId
        LIMIT ? OFFSET ?
    `

	accounts := make([]*models.VendorModel, 0)
	if err := s.db.SelectContext(ctx, &accounts, query, limit, offset); err != nil {
		return nil, err
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

	getSocialsQuery := `
        SELECT
            url
        FROM
            Social
        WHERE
            userId = ?
    `

	for _, account := range accounts {
		if restricted, err := s.IsRestricted(account.VendorId); err != nil {
			return nil, err
		} else {
			account.Restricted = restricted
		}

		expertise := make([]string, 0)
		if err := s.db.SelectContext(ctx, &expertise, expertiseQuery, account.VendorId); err != nil {
			return nil, err
		}
		account.Expertise = expertise

		socials := make([]string, 0)
		if err := s.db.SelectContext(ctx, &socials, getSocialsQuery, account.VendorId); err != nil {
			return nil, err
		}
		account.Socials = socials
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return accounts, nil
}

func (s *MysqlVendorRepository) FindByEmailHash(emailHash string) (*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	vendor := new(models.VendorModel)
	query := `
        SELECT  
            v.joinedAt,
            v.vendorId,
            v.rating,
            u.name AS name,
            u.email AS email,
            u.phone AS phone,
            u.imageUrl AS imageUrl
        FROM 
            Vendor  v
            JOIN User u ON u.id = v.vendorId
        WHERE 
            u.emailHash = ?
    `
	if err := s.db.GetContext(ctx, vendor, query, emailHash); err != nil {
		return nil, err
	}

	if restricted, err := s.IsRestricted(vendor.VendorId); err != nil {
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
	if err := s.db.SelectContext(ctx, &expertise, expertiseQuery, vendor.VendorId); err != nil {
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
	if err := s.db.SelectContext(ctx, &socials, getSocialsQuery, vendor.VendorId); err != nil {
		return nil, err
	}
	vendor.Socials = socials

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}

func (s *MysqlVendorRepository) FindById(id string) (*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	vendor := new(models.VendorModel)
	query := `
        SELECT  
            v.vendorId,
            v.joinedAt,
            v.vendorId,
            v.rating,
            u.name AS name,
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

	if _, date, err := s.IsVerified(id); err != nil {
		return nil, err
	} else {
		vendor.VerifiedAt = date
	}

	if restricted, err := s.IsRestricted(vendor.VendorId); err != nil {
		return nil, err
	} else {
		vendor.Restricted = restricted
	}

	if banned, err := s.IsBanned(vendor.VendorId); err != nil {
		return nil, err
	} else {
		vendor.Banned = banned
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

	getAddressQuery := `
        SELECT
            u.address
        FROM
            User u
            JOIN Vendor v ON u.id = v.vendorId
        WHERE
            v.vendorId = ?
    `
	if err := s.db.GetContext(ctx, &vendor.Address, getAddressQuery, id); err != nil {
		return nil, err
	}

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
            format(rate, 2) as rate,
            latitude, 
            longitude,
            createdAt,
            updatedAt
        FROM 
            Service
        WHERE
            vendorId = ?
    `

	services := make([]*models.ServiceModel, 0)
	if err := s.db.SelectContext(ctx, &services, query, vendorId); err != nil {
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

func (s *MysqlVendorRepository) IsVerified(userId string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT updatedAt FROM IdentityVerification WHERE userId = ? AND status = 'approved' LIMIT 1"
	var updatedAt sql.NullString
	if err := s.db.GetContext(ctx, &updatedAt, query, userId); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return false, "", nil
		}

		return false, "", err
	}

	if !updatedAt.Valid {
		return false, "", nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, "", context.DeadlineExceeded
	}

	return true, updatedAt.String, nil
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

func (s *MysqlVendorRepository) IsBanned(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            CASE
                WHEN EXISTS (SELECT 1 FROM Ban WHERE userId = ?)
                THEN 1
                ELSE 0
            END AS user_exists;
    `

	isBanned := false
	if err := s.db.GetContext(ctx, &isBanned, query, userId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isBanned, nil
}

func (s *MysqlVendorRepository) AddExpertise(userId, expertiseId, supportingImageId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO
            UserExpertise (userId, expertiseId, supportingImage)
        VALUES
            (?, ?, ?)
    `

	if _, err := s.db.ExecContext(ctx, query, userId, expertiseId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlVendorRepository) GetAllExpertise(userId string) ([]*models.UserExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            ue.userId,
            ue.expertiseId,
            ue.supportingImage,
            ue.createdAt,
            e.title AS expertise,
            a.createdAt AS dateApplied,
            a.updatedAt AS dateApproved,
            si.url AS supportingDocumentImage
        FROM
            UserExpertise ue
            JOIN Expertise e ON e.id = ue.expertiseId
            JOIN Application a ON a.applicantId = ue.userId
            JOIN SupportingImage si ON si.id = ue.supportingImage
        WHERE
            ue.userId = ?
        ORDER BY
            ue.createdAt DESC
    `

	expertises := make([]*models.UserExpertiseModel, 0)
	if err := s.db.SelectContext(ctx, &expertises, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return expertises, nil
}

func (s *MysqlVendorRepository) GetBookingsWithStatus(vendorId, status string) ([]*models.BookingModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	allowedStatus := []string{"all", "pending", "confirmed", "done", "rejected", "cancelled"}
	if !slices.Contains(allowedStatus, status) {
		return nil, errors.New("invalid_status")
	}

	query := `
        SELECT
            b.*
        FROM
            Booking b
            JOIN Vendor v ON v.vendorId = b.vendorId
        WHERE
            b.vendorId = ?
    `
	if status != "all" {
		query += fmt.Sprintf(" AND b.status = '%s'", status)
	}

	bookings := make([]*models.BookingModel, 0)
	if err := s.db.SelectContext(ctx, &bookings, query, vendorId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bookings, nil
}
