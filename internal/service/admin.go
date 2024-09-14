package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/store/admin"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AdminService struct {
	store     admin.AdminStore
	encryptor auth.Encryption
	jwt       auth.Authenticator
	hash      auth.Hash
}

func NewAdminService(store admin.AdminStore, encryptor auth.Encryption, jwt auth.Authenticator, hash auth.Hash) *AdminService {
	return &AdminService{
		store:     store,
		encryptor: encryptor,
		jwt:       jwt,
		hash:      hash,
	}
}

func (s *AdminService) BaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, "admin")
}

func (s *AdminService) Login(c echo.Context) error {
	req := new(request.AdminLoginPayload)
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

	usernameHash, err := s.hash.Generate([]byte(req.Username))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error hashing username",
			Error:   err.Error(),
		})
	}

	admin, err := s.store.FindByUsernameHash(usernameHash)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Invalid credentials",
			Error:   err.Error(),
		})
	}

	if auth.IsPasswordMatch(admin.Password, req.Password) == false {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Invalid credentials",
			Error:   "Invalid credentials",
		})
	}

	accessToken, err := s.jwt.GenerateAdminAccessToken(auth.AdminJWTClaims{
		Id:       admin.Id,
		Username: req.Username,
		Role:     admin.Role,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.REFRESH_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	session := models.NewSessionModel(refreshToken)
	if err := s.store.Login(session); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating session",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"role":         admin.Role,
		"adminId":      admin.Id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (s *AdminService) Refresh(c echo.Context) error {
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
	adminId, err := utils.GetAdminIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	admin, err := s.store.FindById(adminId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Admin not found",
			Error:   "Admin not found",
		})
	}

	if plain, err := s.encryptor.DecryptString(admin.Username); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		admin.Username = plain
	}

	accessToken, err := s.jwt.GenerateAdminAccessToken(auth.AdminJWTClaims{
		Id:       adminId,
		Username: admin.Username,
		Role:     admin.Role,
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

func (s *AdminService) Logout(c echo.Context) error {
	req := new(request.LogoutPayload)
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

	if err := s.store.DoesRefreshTokenExists(req.RefreshToken); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Session not found",
			Error:   err.Error(),
		})
	}

	if err := s.store.Logout(req.RefreshToken); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error logging out",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
