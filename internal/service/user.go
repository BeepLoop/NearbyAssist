package service

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/store/user"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserService struct {
	store     user.UserStore
	encryptor auth.Encryption
}

func NewUserService(store user.UserStore, encryptor auth.Encryption) *UserService {
	return &UserService{
		store:     store,
		encryptor: encryptor,
	}
}

func (s *UserService) BaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, "me")
}

func (s *UserService) Login(c echo.Context) error {
	req := new(request.UserLoginPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	emailHash, err := auth.Sha256([]byte(req.Email))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error hashing username",
			Error:   err.Error(),
		})
	}

	existingUser, err := s.store.FindByEmailHash(emailHash)
	if err != nil {
		// If user is not found, continue to registration
		return s.Register(req, emailHash, c)
	}

	accessToken, err := auth.GenerateUserAccessToken(auth.UserJWTClaims{
		Id:    existingUser.Id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.REFRESH_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"user": models.UserModel{
			Model:    models.Model{Id: existingUser.Id},
			Name:     req.Name,
			Email:    req.Email,
			ImageUrl: req.Image,
			Verified: existingUser.Verified,
		},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (s *UserService) Register(req *request.UserLoginPayload, emailHash string, c echo.Context) error {
	newUser := new(models.UserModel)
	newUser.EmailHash = emailHash
	newUser.ImageUrl = req.Image

	if cipher, err := s.encryptor.EncryptString(req.Name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ENCRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		newUser.Name = cipher
	}

	if cipher, err := s.encryptor.EncryptString(req.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ENCRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		newUser.Email = cipher
	}

	if _, err := s.store.CreateUser(newUser); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating user",
			Error:   err.Error(),
		})
	}

	accessToken, err := auth.GenerateUserAccessToken(auth.UserJWTClaims{
		Id:    newUser.Id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.REFRESH_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"user": models.UserModel{
			Model:    models.Model{Id: newUser.Id},
			Name:     req.Name,
			Email:    req.Email,
			ImageUrl: req.Image,
			Verified: newUser.Verified,
		},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (s *UserService) Refresh(c echo.Context) error {
	req := new(request.TokenRefreshPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	// Check if refreshToken exists
	if err := s.store.DoesRefreshTokenExists(req.RefreshToken); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Session not found",
			Error:   err.Error(),
		})
	}

	// Check if refreshToken is blacklisted
	if err := s.store.IsRefreshTokenBlacklisted(req.RefreshToken); err == nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Token blacklisted",
			Error:   "Token blacklisted",
		})
	}

	// Generate new accessToken
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "User Id not found in claims",
			Error:   "User Id not found in claims",
		})
	}

	user, err := s.store.FindById(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "User not found",
			Error:   "User not found",
		})
	}

	if plain, err := s.encryptor.DecryptString(user.Name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		user.Name = plain
	}

	if plain, err := s.encryptor.DecryptString(user.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		user.Email = plain
	}

	accessToken, err := auth.GenerateUserAccessToken(auth.UserJWTClaims{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"accessToken": accessToken,
	})
}
