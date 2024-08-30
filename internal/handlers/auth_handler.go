package handlers

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/hash"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authHandler struct {
	server *server.Server
}

func NewAuthHandler(server *server.Server) *authHandler {
	return &authHandler{server}
}

func (h *authHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Auth route is up and running!",
	})
}

func (h *authHandler) HandleAdminLogin(c echo.Context) error {
	// Bind request body
	req := new(request.AdminLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate required fields
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	admin := models.NewAdminModel(h.server.IdGen, h.server.DB)
	if admin == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if usernameHash, err := h.server.Hash.Hash([]byte(req.Username)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		admin.UsernameHash = usernameHash
	}

	if _, err := admin.FindByUsernameHash(); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	if admin.IsPasswordMatch(req.Password) != true {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	accessToken, err := h.server.Auth.GenerateAdminAccessToken(authenticator.AdminOptions{
		Id:       admin.Id,
		Username: req.Username,
		Role:     admin.Role,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, authenticator.ACCESS_TOKEN_ERR)
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, authenticator.REFRESH_TOKEN_ERR)
	}

	session := models.NewSessionModel(refreshToken, h.server.IdGen, h.server.DB)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}
	if _, err := session.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating session")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"role":         admin.Role,
		"adminId":      admin.Id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *authHandler) HandleClientLogin(c echo.Context) error {
	req := new(request.UserLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := models.NewUserModel(h.server.IdGen, h.server.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if emailHash, err := h.server.Hash.Hash([]byte(req.Email)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		user.Hash = emailHash
	}

	if _, err := user.FindByEmailHash(); err != nil {
		// Account not found

		if encrypted, err := h.server.Encrypt.EncryptString(req.Name); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
		} else {
			user.Name = encrypted
		}

		if encrypted, err := h.server.Encrypt.EncryptString(req.Email); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
		} else {
			user.Email = encrypted
		}

		if _, err := user.Create(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error creating account")
		}
	}

	accessToken, err := h.server.Auth.GenerateUserAccessToken(authenticator.UserOptions{
		Id:    user.Id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, authenticator.ACCESS_TOKEN_ERR)
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, authenticator.REFRESH_TOKEN_ERR)
	}

	session := models.NewSessionModel(refreshToken, h.server.IdGen, h.server.DB)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}
	if _, err := session.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating session")
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"user": models.UserModel{
			Model:    models.Model{Id: user.Id},
			Name:     req.Name,
			Email:    req.Email,
			ImageUrl: req.Image,
			Verified: user.IsVerified(),
		},
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *authHandler) HandleLogout(c echo.Context) error {
	req := new(request.Logout)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session := models.NewSessionModel(req.Token, h.server.IdGen, h.server.DB)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := session.GetIfActive(); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "No active session found")
	}

	if err := session.Logout(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error logging out session")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Logout successful",
	})
}

func (h *authHandler) HandleTokenRefresh(c echo.Context) error {
	req := new(request.RefreshToken)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session := models.NewSessionModel(req.Token, h.server.IdGen, h.server.DB)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := session.FindByToken(); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid token")
	}

	if blacklisted, err := session.IsBlacklisted(); err != nil || blacklisted {
		return echo.NewHTTPError(http.StatusForbidden, "Token is blacklisted")
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	var newAccessToken string
	if _, ok := claims["role"].(string); ok {
		// Admin requests refresh token
		idFromJWT, ok := claims["adminId"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Admin Id not found in claims")
		}

		admin := models.NewAdminModel(h.server.IdGen, h.server.DB)
		if admin == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
		}

		if _, err := admin.FindById(idFromJWT); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "ID from JWT not found")
		}

		if _, err := admin.DecryptUsername(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
		}

		accessToken, err := h.server.Auth.GenerateAdminAccessToken(authenticator.AdminOptions{
			Id:       admin.Id,
			Username: admin.Username,
			Role:     admin.Role,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, authenticator.ACCESS_TOKEN_ERR)
		}

		newAccessToken = accessToken
	} else {
		// User requests refresh token
		idFromJWT, ok := claims["userId"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Admin Id not found in claims")
		}

		user := models.NewUserModel(h.server.IdGen, h.server.DB)
		if user == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
		}

		if _, err := user.FindById(idFromJWT); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "ID from JWT not found")
		}

		if _, err := user.DecryptName(h.server.Encrypt.EncryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
		}
		if _, err := user.DecryptEmail(h.server.Encrypt.EncryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
		}

		accessToken, err := h.server.Auth.GenerateUserAccessToken(authenticator.UserOptions{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, authenticator.ACCESS_TOKEN_ERR)
		}

		newAccessToken = accessToken
	}

	// As protection against my dumb self, check first if newAccessToken is an empty string
	if newAccessToken == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while determining if request is from user or admin")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"accessToken": newAccessToken,
	})
}
