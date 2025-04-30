package user_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlUserRepository struct {
	db *sqlx.DB
}

func NewMysqlUserRepository(db *sqlx.DB) *MysqlUserRepository {
	return &MysqlUserRepository{db: db}
}

func (s *MysqlUserRepository) CreateUser(user *models.UserModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	user.Id = utils.GenerateUserId()
	createUser := `
        INSERT INTO
            User (id, name, email, imageUrl, emailHash, phone)
        VALUES
            (:id, :name, :email, :imageUrl, :emailHash, :phone)
    `
	if _, err := tx.NamedExecContext(ctx, createUser, user); err != nil {
		return "", err
	}

	createAddress := `
        INSERT INTO
            Address (id, address, latitude, longitude)
        VALUES
            (:id, :address, :latitude, :longitude)
    `
	user.Address.Id = utils.GenerateId()
	if _, err := tx.NamedExecContext(ctx, createAddress, user.Address); err != nil {
		return "", err
	}

	createAddressRelation := `
        INSERT INTO
            UserAddress (userId, addressId)
        VALUES
            (?, ?)
    `
	if _, err := tx.ExecContext(ctx, createAddressRelation, user.Id, user.Address.Id); err != nil {
		return "", err
	}

	createIdentification := `
        INSERT INTO
            Identification (id, type, referenceNumber, frontImageUrl, backImageUrl, selfieImageUrl)
        VALUES
            (:id, :type, :referenceNumber, :frontImageUrl, :backImageUrl, :selfieImageUrl)
    `
	user.Identification.Id = utils.GenerateId()
	if _, err := tx.NamedExecContext(ctx, createIdentification, user.Identification); err != nil {
		return "", err
	}

	createIdentificationRelation := `
        INSERT INTO
            UserIdentification (userId, identificationId)
        VALUES
            (?, ?)
    `
	if _, err := tx.ExecContext(ctx, createIdentificationRelation, user.Id, user.Identification.Id); err != nil {
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

	return user.Id, nil
}

func (s *MysqlUserRepository) GetBasicUserAccounts(limit, offset int) ([]*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getAccountsQuery := `
        SELECT
            u.id
        FROM
            User u
            LEFT JOIN Vendor v ON v.vendorId = u.id
        WHERE
            v.vendorId IS NULL
        ORDER BY
            createdAt DESC
        LIMIT ? OFFSET ?
    `
	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, getAccountsQuery, limit, offset); err != nil {
		return nil, err
	}

	accounts := make([]*models.UserModel, 0)
	for _, id := range ids {
		account, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return accounts, nil
}

func (s *MysqlUserRepository) GetAllUserAccounts(limit, offset int) ([]*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            id
        FROM
            User
        ORDER BY
            createdAt DESC
        LIMIT ? OFFSET ?
    `
	ids := make([]string, 0)
	if err := s.db.SelectContext(ctx, &ids, query, limit, offset); err != nil {
		return nil, err
	}

	accounts := make([]*models.UserModel, 0)
	for _, id := range ids {
		account, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return accounts, nil
}

func (s *MysqlUserRepository) FindById(id string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	getUserQuery := `
        SELECT 
            id, name, email, imageUrl, phone, verified, verifiedAt, createdAt
        FROM 
            User 
        WHERE 
            id = ?
    `

	user := new(models.UserModel)
	err := s.db.GetContext(ctx, user, getUserQuery, id)
	if err != nil {
		return nil, err
	}

	if banned, err := s.IsBanned(user.Id); err != nil {
		return nil, err
	} else {
		user.Banned = banned
	}

	if restricted, expired, err := s.IsRestricted(user.Id); err != nil {
		return nil, err
	} else {
		user.Restricted = restricted && !expired
	}

	if address, err := s.GetAddress(user.Id); err != nil {
		return nil, err
	} else {
		user.Address = *address
	}

	if identification, err := s.GetIdentification(user.Id); err != nil {
		return nil, err
	} else {
		user.Identification = *identification
	}

	if socials, err := s.GetSocials(user.Id); err != nil {
		return nil, err
	} else {
		user.Socials = slices.AppendSeq(
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

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (s *MysqlUserRepository) FindByEmailHash(emailHash string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user := new(models.UserModel)

	query := `
        SELECT 
            id, name, email, imageUrl, phone, verified, verifiedAt, createdAt
        FROM 
            User 
        WHERE 
            emailHash = ?
    `
	if err := s.db.GetContext(ctx, user, query, emailHash); err != nil {
		return nil, err
	}

	if banned, err := s.IsBanned(user.Id); err != nil {
		return nil, err
	} else {
		user.Banned = banned
	}

	if restricted, expired, err := s.IsRestricted(user.Id); err != nil {
		return nil, err
	} else {
		user.Restricted = restricted && !expired
	}

	if address, err := s.GetAddress(user.Id); err != nil {
		return nil, err
	} else {
		user.Address = *address
	}

	if identification, err := s.GetIdentification(user.Id); err != nil {
		return nil, err
	} else {
		user.Identification = *identification
	}

	if socials, err := s.GetSocials(user.Id); err != nil {
		return nil, err
	} else {
		user.Socials = slices.AppendSeq(
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

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (s *MysqlUserRepository) GetAddress(userId string) (*models.AddressModel, error) {
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

func (s *MysqlUserRepository) GetSocials(userId string) ([]*models.SocialModel, error) {
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

func (s *MysqlUserRepository) GetIdentification(userId string) (*models.IdentificationModel, error) {
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

func (s *MysqlUserRepository) IsVendor(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	isVendor := false
	query := "SELECT EXISTS (SELECT 1 FROM Vendor WHERE vendorId = ?) AS is_vendor"
	if err := s.db.GetContext(ctx, &isVendor, query, userId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isVendor, nil
}

func (s *MysqlUserRepository) GetExpertise(userId string) ([]*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get vendor expertise
	expertiseQuery := `
        SELECT 
            e.id,
            e.title
        FROM 
            Expertise e
            JOIN UserExpertise ve ON ve.expertiseId = e.id
        WHERE
            ve.userId = ?
    `

	expertises := make([]*models.ExpertiseModel, 0)
	if err := s.db.SelectContext(ctx, &expertises, expertiseQuery, userId); err != nil {
		return nil, err
	}

	// Get all tags for each expertise
	tagQuery := `
        SELECT
            t.id,
            t.title
        FROM 
            Tag t 
            JOIN ExpertiseTag et ON et.tagId = t.id
        WHERE
            et.expertiseId = ?
    `

	for _, expertise := range expertises {
		tags := make([]*models.TagModel, 0)
		if err := s.db.SelectContext(ctx, &tags, tagQuery, expertise.Id); err != nil {
			return nil, err
		}

		expertise.Tags = tags
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return expertises, nil
}

func (s *MysqlUserRepository) AddSocial(data *models.SocialModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addSocialQuery := `
        INSERT INTO
            Social (id, userId, site, title, url)
        VALUES
            (:id, :userId, :site, :title, :url)
    `

	data.Id = utils.GenerateId()
	if _, err := s.db.NamedExecContext(ctx, addSocialQuery, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlUserRepository) DeleteSocial(userId, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deleteSocialQuery := "DELETE FROM Social WHERE userId = ? AND id = ?"
	if _, err := s.db.ExecContext(ctx, deleteSocialQuery, userId, id); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlUserRepository) IsBanned(userId string) (bool, error) {
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

func (s *MysqlUserRepository) BanUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	banQuery := "INSERT INTO Ban (userId) VALUES(?)"
	if _, err := s.db.ExecContext(ctx, banQuery, userId); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil
		}

		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlUserRepository) UnbanUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	unbanQuery := "DELETE FROM Ban WHERE userId = ?"
	if _, err := s.db.ExecContext(ctx, unbanQuery, userId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

// Return isRestricted, isExpired, error
func (s *MysqlUserRepository) IsRestricted(userId string) (bool, bool, error) {
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

func (s *MysqlUserRepository) RestrictUser(data *models.RestrictionModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO
            Restricted (userId, reason, endTime)
        VALUES
            (:userId, :reason, :endTime)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil
		}

		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlUserRepository) LiftRestrictionIfExpired(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM Restricted WHERE userId = ? AND endTime <= ?"
	if _, err := s.db.ExecContext(ctx, query, userId, utils.CurrentTimeStamp()); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlUserRepository) ForceLiftRestriction(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM Restricted WHERE userId = ?"
	if _, err := s.db.ExecContext(ctx, query, userId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
