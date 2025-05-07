package vendor_repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
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

func (s *MysqlVendorRepository) FindByEmailHash(emailHash string) (*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	vendor := new(models.VendorModel)
	findVendorQuery := `
        SELECT  
            v.vendorId, v.dbl, v.joinedAt, v.rating
        FROM 
            Vendor  v
            JOIN User u ON u.id = v.vendorId
        WHERE 
            u.emailHash = ?
    `
	if err := s.db.GetContext(ctx, vendor, findVendorQuery, emailHash); err != nil {
		return nil, err
	}

	getUserQuery := `
        SELECT 
            id, name, email, imageUrl, phone, verified, verifiedAt, createdAt
        FROM 
            User 
        WHERE 
            id = ?
    `

	if err := s.db.GetContext(ctx, &vendor.User, getUserQuery, vendor.VendorId); err != nil {
		return nil, err
	}

	if identification, err := s.getIdentification(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Identification = *identification
	}

	if banned, err := s.isBanned(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Banned = banned
	}

	if restricted, expired, err := s.IsRestricted(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Restricted = restricted && !expired
	}

	if address, err := s.getAddress(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Address = *address
	}

	if socials, err := s.getSocials(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Socials = slices.AppendSeq(
			make([]models.SocialModel, 0),
			utils.Map(socials, func(social *models.SocialModel) models.SocialModel {
				return models.SocialModel{
					Model:  social.Model,
					UserId: social.UserId,
					Site:   social.Site,
					Title:  social.Title,
					Url:    social.Url,
				}
			}),
		)
	}

	if list, err := s.getVendorExpertise(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.Expertise = slices.AppendSeq(
			make([]models.ExpertiseModel, 0),
			utils.Map(list, func(e *models.ExpertiseModel) models.ExpertiseModel {
				return models.ExpertiseModel{
					Model: models.Model{Id: e.Id, CreatedAt: e.CreatedAt},
					Title: e.Title,
				}
			}),
		)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}

func (s *MysqlVendorRepository) FindById(vendorId string) (*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	vendor := new(models.VendorModel)
	findVendorQuery := `
        SELECT  
            v.vendorId, v.dbl, v.joinedAt, v.rating
        FROM 
            Vendor  v
            JOIN User u ON u.id = v.vendorId
        WHERE 
            u.id = ?
    `
	if err := s.db.GetContext(ctx, vendor, findVendorQuery, vendorId); err != nil {
		return nil, err
	}

	getUserQuery := `
        SELECT 
            id, name, email, imageUrl, phone, verified, verifiedAt, createdAt
        FROM 
            User 
        WHERE 
            id = ?
    `

	if err := s.db.GetContext(ctx, &vendor.User, getUserQuery, vendor.VendorId); err != nil {
		return nil, err
	}

	if identification, err := s.getIdentification(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Identification = *identification
	}

	if banned, err := s.isBanned(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Banned = banned
	}

	if restricted, expired, err := s.IsRestricted(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Restricted = restricted && !expired
	}

	if address, err := s.getAddress(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Address = *address
	}

	if socials, err := s.getSocials(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.User.Socials = slices.AppendSeq(
			make([]models.SocialModel, 0),
			utils.Map(socials, func(social *models.SocialModel) models.SocialModel {
				return models.SocialModel{
					Model:  social.Model,
					UserId: social.UserId,
					Site:   social.Site,
					Title:  social.Title,
					Url:    social.Url,
				}
			}),
		)
	}

	if list, err := s.getVendorExpertise(vendor.User.Id); err != nil {
		return nil, err
	} else {
		vendor.Expertise = slices.AppendSeq(
			make([]models.ExpertiseModel, 0),
			utils.Map(list, func(e *models.ExpertiseModel) models.ExpertiseModel {
				return models.ExpertiseModel{
					Model:              models.Model{Id: e.Id, CreatedAt: e.CreatedAt},
					UpdateableModel:    e.UpdateableModel,
					Title:              e.Title,
					Tags:               e.Tags,
					DateApplied:        e.DateApplied,
					DateApproved:       e.DateApproved,
					SupportingImageUrl: e.SupportingImageUrl,
				}
			}),
		)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}

func (s *MysqlVendorRepository) GetAll(limit, offset int) ([]*models.VendorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT  
            v.vendorId, v.dbl, v.rating, v.joinedAt
        FROM 
            Vendor  v
            JOIN User u ON u.id = v.vendorId
        LIMIT ? OFFSET ?
    `

	vendors := make([]*models.VendorModel, 0)
	if err := s.db.SelectContext(ctx, &vendors, query, limit, offset); err != nil {
		return nil, err
	}

	for _, vendor := range vendors {
		if res, err := s.FindById(vendor.VendorId); err != nil {
			return nil, err
		} else {
			vendor.User = res.User
			vendor.Expertise = res.Expertise
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendors, nil
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
            FORMAT(rate, 2) AS rate,
            createdAt,
            disabled
        FROM 
            Service
        WHERE
            vendorId = ?
        ORDER BY
            createdAt, updatedAt DESC
    `
	services := make([]*models.ServiceModel, 0)
	if err := s.db.SelectContext(ctx, &services, query, vendorId); err != nil {
		return nil, err
	}

	for _, service := range services {
		if extras, err := s.getExtras(service.Id); err != nil {
			return nil, err
		} else {
			service.Extras = extras
		}

		if tags, err := s.getTags(service.Id); err != nil {
			return nil, err
		} else {
			service.Tags = tags
		}

		if images, err := s.getPhotos(service.Id); err != nil {
			return nil, err
		} else {
			service.Images = images
		}

		if address, err := s.getAddress(service.VendorId); err != nil {
			return nil, err
		} else {
			service.Address = *address
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}

func (s *MysqlVendorRepository) getTags(serviceId string) ([]*models.TagModel, error) {
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

func (s *MysqlVendorRepository) getExtras(serviceId string) ([]*models.ExtraModel, error) {
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

func (s *MysqlVendorRepository) getPhotos(serviceId string) ([]*models.ServicePhotoModel, error) {
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

// Return isRestricted, isExpired, error
func (s *MysqlVendorRepository) IsRestricted(userId string) (bool, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isRestrictedQuery := `
        SELECT
            CASE
                WHEN (SELECT 1 FROM Restricted WHERE userId = ?)
                THEN 1
                ELSE 0
            END AS user_exists;
    `
	isRestricted := false
	if err := s.db.GetContext(ctx, &isRestricted, isRestrictedQuery, userId); err != nil {
		return false, false, err
	}

	if !isRestricted {
		return false, false, nil
	}

	isRestrictionExpiredQuery := `
        SELECT
            CASE
                WHEN (SELECT 1 FROM Restricted WHERE userId = ? AND endTime < ?)
                THEN 1
                ELSE 0
            END AS isExpired;
    `
	isExpired := false
	if err := s.db.GetContext(ctx, &isExpired, isRestrictionExpiredQuery, userId, utils.CurrentTimeStamp()); err != nil {
		return false, false, nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, false, context.DeadlineExceeded
	}

	return isRestricted, isExpired, nil
}

func (s *MysqlVendorRepository) isBanned(userId string) (bool, error) {
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

	if _, err := s.db.ExecContext(ctx, query, userId, expertiseId, supportingImageId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
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

func (s *MysqlVendorRepository) getAddress(userId string) (*models.AddressModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getAddress := `
        SELECT
            a.id, a.address, a.latitude, a.longitude
        FROM
            Address a
            JOIN UserAddress ua ON ua.addressId = a.id
        WHERE
            ua.userId = ?
    `

	address := new(models.AddressModel)
	if err := s.db.GetContext(ctx, address, getAddress, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return address, nil
}

func (s *MysqlVendorRepository) getSocials(userId string) ([]*models.SocialModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            id, userId, site, title, url, createdAt
        FROM
            Social
        WHERE
            userId = ?
    `

	socials := make([]*models.SocialModel, 0)
	if err := s.db.SelectContext(ctx, &socials, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return socials, nil
}

func (s *MysqlVendorRepository) getVendorExpertise(vendorId string) ([]*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            e.id,
            e.title,
            e.createdAt,
            a.createdAt AS dateApplied,
            a.updatedAt AS dateApproved,
            i.url AS supportingImageUrl
        FROM
            Expertise e
            JOIN UserExpertise ue ON ue.expertiseId = e.id
            JOIN Application a ON a.applicantid = ue.userId
            JOIN SupportingImage i ON i.id = ue.supportingImage
        WHERE
            ue.userId = ? AND a.status = 'approved'
    `

	expertise := make([]*models.ExpertiseModel, 0)
	if err := s.db.SelectContext(ctx, &expertise, query, vendorId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return expertise, nil
}

func (s *MysqlVendorRepository) getIdentification(userId string) (*models.IdentificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            i.type,
            i.referenceNumber,
            i.frontImageUrl,
            i.backImageUrl,
            i.selfieImageUrl
        FROM
            Identification i
            JOIN UserIdentification ui ON ui.identificationId = i.id
        WHERE
            ui.userId = ?
    `

	identification := new(models.IdentificationModel)
	if err := s.db.GetContext(ctx, identification, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return identification, nil
}

func (s *MysqlVendorRepository) CompletedBookingCountOfService(vendorId, serviceId string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            COUNT(id)
        FROM
            Booking
        WHERE
            vendorId = ? AND status = 'done' AND serviceId = ?
    `
	count := 0
	if err := s.db.GetContext(ctx, &count, query, vendorId, serviceId); err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlVendorRepository) HasExpertise(vendorId, expertiseId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT EXISTS
            (SELECT 1 FROM UserExpertise WHERE userId = ? AND expertiseId = ?)
        AS has_expertise
    `
	hasExpertise := false
	if err := s.db.GetContext(ctx, &hasExpertise, query, vendorId, expertiseId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return hasExpertise, nil
}

func (s *MysqlVendorRepository) GetPoliceClearance(vendorId string) (*models.PoliceClearanceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            p.*
        FROM
            Vendor v
            JOIN Application a ON a.applicantId = v.vendorId
            JOIN PoliceClearance p ON p.id = a.policeClearance
        WHERE
            v.vendorId = ?
        LIMIT 1
    `
	clearance := new(models.PoliceClearanceModel)
	if err := s.db.GetContext(ctx, clearance, query, vendorId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return clearance, nil
}

func (s *MysqlVendorRepository) IsFullyBookedAt(vendorId, schedule string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT CASE
            WHEN (
                SELECT COUNT(*)
                FROM Booking b
                WHERE b.vendorId = ? AND DATE(b.scheduledAt) = DATE(?)
            ) >= (
                SELECT v.dbl
                FROM Vendor v
                WHERE v.vendorId = ?
            )
            THEN 1
            ELSE 0
        END AS fully_booked
    `

	fullyBooked := false
	if err := s.db.GetContext(ctx, &fullyBooked, query, vendorId, schedule, vendorId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return fullyBooked, nil
}

func (s *MysqlVendorRepository) SetDBL(vendorId string, value int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Vendor SET dbl = ? WHERE vendorId = ?"
	if _, err := s.db.ExecContext(ctx, query, value, vendorId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
