package user_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}
	user.Id = id

	query := "INSERT INTO User (id, name, email, imageUrl, emailHash) VALUES (:id, :name, :email, :imageUrl, :emailHash)"
	if _, err := s.db.NamedExecContext(ctx, query, user); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return id, nil
}

func (s *MysqlUserRepository) Login(data *models.SessionModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := gonanoid.New(); err != nil {
		return err
	} else {
		data.Id = id
	}

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

	id, err := gonanoid.New()
	if err != nil {
		return err
	}

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
            banned,
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

	getUserQuery := "SELECT id, name, email, imageUrl, address, banned, createdAt FROM User where id = ?"
	user := new(models.UserModel)
	if err := s.db.GetContext(ctx, user, getUserQuery, userId); err != nil {
		return nil, err
	}
	accountData.Id = user.Id
	accountData.ProfileURL = user.ImageUrl
	accountData.Name = user.Name
	accountData.Email = user.Email
	accountData.Address = user.Address
	accountData.Banned = user.Banned
	accountData.CreatedAt = user.CreatedAt

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

	getServicesQuery := "SELECT * FROM Service WHERE vendorId = ?"
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

	query := "SELECT id, name, email, imageUrl, verified, banned, address, phone, latitude, longitude FROM User WHERE emailHash = ?"
	if err := s.db.GetContext(ctx, user, query, emailHash); err != nil {
		return nil, err
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

	if generatedId, err := gonanoid.New(); err != nil {
		return err
	} else {
		data.Id = generatedId
	}

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

func (s *MysqlUserRepository) BanUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE User SET banned = 1 WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, userId); err != nil {
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

	query := "UPDATE User SET banned = 0 WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, userId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
