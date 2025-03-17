package user_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
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

	user.Id = utils.GenerateUserId()

	query := "INSERT INTO User (id, name, email, imageUrl, emailHash) VALUES (:id, :name, :email, :imageUrl, :emailHash)"
	if _, err := s.db.NamedExecContext(ctx, query, user); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return user.Id, nil
}

func (s *MysqlUserRepository) Login(data *models.SessionModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	query := "INSERT INTO Session (id, refreshToken) VALUES (:id, :refreshToken)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlUserRepository) Logout(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	updateSession := "UPDATE Session SET status = 'offline' WHERE refreshToken = ? AND status = 'online'"
	if _, err := s.db.ExecContext(ctx, updateSession, refreshToken); err != nil {
		return err
	}

	id := utils.GenerateId()
	blacklistToken := `INSERT INTO Blacklist (id, token) VALUES (?, ?)`
	if _, err := tx.ExecContext(ctx, blacklistToken, id, refreshToken); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlUserRepository) GetBasicUserAccounts(limit, offset int) ([]*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	accounts := make([]*models.UserModel, 0)

	getAccountsQuery := `
        SELECT
            u.id, u.name, u.email, u.imageUrl, u.verified, u.createdAt
        FROM
            User u
            LEFT JOIN Vendor v ON v.vendorId = u.id
        WHERE
            v.vendorId IS NULL
        ORDER BY createdAt DESC
        LIMIT ? OFFSET ?
    `
	if err := tx.SelectContext(ctx, &accounts, getAccountsQuery, limit, offset); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}

		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return accounts, nil
}

func (s *MysqlUserRepository) GetAllUserAccounts(limit, offset int) ([]*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	accounts := make([]*models.UserModel, 0)

	getAccountsQuery := `
        SELECT
            id, name, email, imageUrl, verified, createdAt
        FROM
            User
        ORDER BY createdAt DESC
        LIMIT ?
        OFFSET ?
    `
	if err := tx.SelectContext(ctx, &accounts, getAccountsQuery, limit, offset); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}

		return nil, err
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
            id,
            name,
            email,
            imageUrl,
            verified,
            address,
            phone,
            latitude,
            longitude,
            createdAt 
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

	getUserSocialsQuery := `SELECT url FROM Social WHERE userId = ?`

	socials := make([]string, 0)
	if err := s.db.SelectContext(ctx, &socials, getUserSocialsQuery, id); err != nil {
		return nil, err
	}
	user.Socials = socials

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (s *MysqlUserRepository) GetUserAccountPageData(userId string) (*models.UserAccountPageData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accountData := new(models.UserAccountPageData)

	getUserQuery := `
        SELECT
            id,
            name,
            email,
            imageUrl,
            verified,
            address,
            createdAt
        FROM 
            User 
        WHERE 
            id = ?
    `
	user := new(models.UserModel)
	if err := s.db.GetContext(ctx, user, getUserQuery, userId); err != nil {
		return nil, err
	}
	accountData.Id = user.Id
	accountData.ProfileURL = user.ImageUrl
	accountData.Name = user.Name
	accountData.Email = user.Email
	accountData.Address = user.Address
	accountData.CreatedAt = user.CreatedAt
	accountData.Verified = user.Verified

	if banned, err := s.IsBanned(user.Id); err != nil {
		return nil, err
	} else {
		accountData.Banned = banned
	}

	if restricted, expired, err := s.IsRestricted(user.Id); err != nil {
		return nil, err
	} else {
		accountData.Restricted = restricted && !expired
	}

	getExpertiseQuery := `
        SELECT
            e.title
        FROM
            VendorExpertise ve
            JOIN Expertise e ON e.id = ve.expertiseId
        WHERE
            ve.vendorId = ?
    `
	expertise := make([]*models.ExpertiseModel, 0)
	if err := s.db.SelectContext(ctx, &expertise, getExpertiseQuery, userId); err != nil {
		return nil, err
	}

	for _, expertise := range expertise {
		accountData.Expertise = append(accountData.Expertise, expertise.Title)
	}

	getServicesQuery := "SELECT * FROM Service WHERE vendorId = ? ORDER BY updatedAT DESC"
	services := make([]*models.ServiceModel, 0)
	if err := s.db.SelectContext(ctx, &services, getServicesQuery, userId); err != nil {
		return nil, err
	}
	accountData.Services = services

	getServiceTagsQuery := `
        SELECT
            t.title
        FROM
            ServiceTag st
            JOIN Tag t ON st.tagId = t.id
        WHERE
            st.serviceId = ?
    `
	for _, service := range accountData.Services {
		tags := make([]string, 0)
		if err := s.db.SelectContext(ctx, &tags, getServiceTagsQuery, service.Id); err != nil {
			return nil, err
		}

		service.TagsAsString = tags
	}

	getServiceImagesQuery := `
        SELECT
            id,
            serviceId,
            vendorId,
            url
        FROM
            ServicePhoto
        WHERE
            serviceId = ?
    `
	for _, service := range accountData.Services {
		images := make([]*models.ServicePhotoModel, 0)
		if err := s.db.SelectContext(ctx, &images, getServiceImagesQuery, service.Id); err != nil {
			return nil, err
		}

		service.Images = images
	}

	getServiceExtrasQuery := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price,
            e.createdAt
        FROM
            ServiceExtra se
            JOIN Extra e ON e.id = se.extraId
        WHERE
            se.serviceId = ?
    `
	for _, service := range accountData.Services {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getServiceExtrasQuery, service.Id); err != nil {
			return nil, err
		}

		service.Extras = extras
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return accountData, nil
}

func (s *MysqlUserRepository) GetSentTransactionCount(userId string) (*models.SentStat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stat := new(models.SentStat)

	query := `
        SELECT 
            SUM(CASE 
                    WHEN YEAR(createdAt) = YEAR(CURDATE()) 
                    AND MONTH(createdAt) = MONTH(CURDATE()) 
                    THEN 1 
                    ELSE 0 
                END) AS currentMonth,
            SUM(CASE 
                    WHEN YEAR(createdAt) = YEAR(CURDATE() - INTERVAL 1 MONTH) 
                    AND MONTH(createdAt) = MONTH(CURDATE() - INTERVAL 1 MONTH) 
                    THEN 1 
                    ELSE 0 
                END) AS lastMonth
        FROM 
            Transaction
        WHERE 
            clientId = ?
    `
	if err := s.db.GetContext(ctx, stat, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return stat, nil
}

func (s *MysqlUserRepository) GetReceivedTransactionCount(userId string) (*models.ReceivedStat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stat := new(models.ReceivedStat)

	query := `
        SELECT 
            SUM(CASE 
                    WHEN YEAR(createdAt) = YEAR(CURDATE()) 
                    AND MONTH(createdAt) = MONTH(CURDATE()) 
                    THEN 1 
                    ELSE 0 
                END) AS currentMonth,
            SUM(CASE 
                    WHEN YEAR(createdAt) = YEAR(CURDATE() - INTERVAL 1 MONTH) 
                    AND MONTH(createdAt) = MONTH(CURDATE() - INTERVAL 1 MONTH) 
                    THEN 1 
                    ELSE 0 
                END) AS lastMonth
        FROM 
            Transaction
        WHERE 
            vendorId = ?
    `
	if err := s.db.GetContext(ctx, stat, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return stat, nil
}

func (s *MysqlUserRepository) FindByEmailHash(emailHash string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user := new(models.UserModel)

	query := `
        SELECT 
            id,
            name,
            email,
            imageUrl,
            verified,
            address,
            phone,
            latitude,
            longitude,
            createdAt
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

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (s *MysqlUserRepository) FindSessionByToken(refreshToken string) (*models.SessionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	session := new(models.SessionModel)
	query := "SELECT id, status, refreshToken FROM Session WHERE refreshToken = ? AND status = 'online'"
	if err := s.db.GetContext(ctx, session, query, refreshToken); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return session, nil
}

func (s *MysqlUserRepository) IsRefreshTokenBlacklisted(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count := 0
	query := "SELECT COUNT(id) FROM Blacklist WHERE token = ?"
	if err := s.db.GetContext(ctx, &count, query, refreshToken); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	if count > 0 {
		return nil
	}

	return errors.New("refreshToken blacklisted")
}

func (s *MysqlUserRepository) IsVendor(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count := 0
	query := "SELECT COUNT(id) FROM Vendor WHERE vendorId = ?"
	if err := s.db.GetContext(ctx, &count, query, userId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
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
            JOIN VendorExpertise ve ON ve.expertiseId = e.id
        WHERE
            ve.vendorId = ?
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

func (s *MysqlUserRepository) AddSocial(data *models.SocialModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addSocialQuery := `
        INSERT INTO
            Social (id, userId, url)
        VALUES
            (:id, :userId, :url)
    `

	data.Id = utils.GenerateId()

	if _, err := s.db.NamedExecContext(ctx, addSocialQuery, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlUserRepository) GetSocials(userId string) ([]*models.SocialModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getSocialsQuery := "SELECT * FROM Social WHERE userId = ?"

	socials := make([]*models.SocialModel, 0)
	if err := s.db.SelectContext(ctx, &socials, getSocialsQuery, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return socials, nil
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
