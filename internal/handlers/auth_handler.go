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
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	// Validate required fields
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	admin := models.NewAdminModel(h.server.IdGen, h.server.DB)
	if admin == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := admin.HashUsername(req.Username, h.server.Hash.Hash); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error hashing username",
			Error:   err.Error(),
		})
	}

	if _, err := admin.FindByUsernameHash(); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Invalid credentials",
			Error:   err.Error(),
		})
	}

	if admin.IsPasswordMatch(req.Password) != true {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Invalid credentials",
			Error:   "Invalid credentials",
		})
	}

	accessToken, err := h.server.Auth.GenerateAdminAccessToken(authenticator.AdminOptions{
		Id:       admin.Id,
		Username: req.Username,
		Role:     admin.Role,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: authenticator.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: authenticator.REFRESH_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	session := models.NewSessionModel(refreshToken, h.server.IdGen, h.server.DB)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}
	if _, err := session.Create(); err != nil {
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

func (h *authHandler) HandleClientLogin(c echo.Context) error {
	req := new(request.UserLogin)
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

	user := models.NewUserModel(h.server.IdGen, h.server.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := user.HashEmail(req.Email, h.server.Hash.Hash); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: hash.HASH_ERROR,
			Error:   err.Error(),
		})
	}

	if _, err := user.FindByEmailHash(); err != nil {
		// Account not found

		if encrypted, err := h.server.Encrypt.EncryptString(req.Name); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: encryption.ENCRYPTION_ERR,
				Error:   err.Error(),
			})
		} else {
			user.Name = encrypted
		}

		if encrypted, err := h.server.Encrypt.EncryptString(req.Email); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: encryption.ENCRYPTION_ERR,
				Error:   err.Error(),
			})
		} else {
			user.Email = encrypted
		}

		if _, err := user.Create(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error creating user",
				Error:   err.Error(),
			})
		}
	}

	accessToken, err := h.server.Auth.GenerateUserAccessToken(authenticator.UserOptions{
		Id:    user.Id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: authenticator.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: authenticator.REFRESH_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	session := models.NewSessionModel(refreshToken, h.server.IdGen, h.server.DB)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}
	if _, err := session.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating session",
			Error:   err.Error(),
		})
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

	session := models.NewSessionModel(req.Token, h.server.IdGen, h.server.DB)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := session.GetIfActive(); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Session not found",
			Error:   err.Error(),
		})
	}

	if err := session.Logout(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error logging out",
			Error:   err.Error(),
		})
	}

	blacklist := models.NewBlacklistModelWithToken(session.RefreshToken, h.server.IdGen, h.server.DB)
	if blacklist == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := blacklist.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error blacklisting token",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Logout successful",
	})
}

func (h *authHandler) HandleTokenRefresh(c echo.Context) error {
	req := new(request.RefreshToken)
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

	session := models.NewSessionModel(req.Token, h.server.IdGen, h.server.DB)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := session.FindByToken(); err != nil {
		println("Error finding session by token")
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Session not found",
			Error:   err.Error(),
		})
	}

	blacklist := models.NewBlacklistModel(h.server.IdGen, h.server.DB)
	if _, err := blacklist.FindByToken(req.Token); err == nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Session token is blacklisted",
			Error:   err.Error(),
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	var newAccessToken string
	if _, ok := claims["role"].(string); ok {
		// Admin requests refresh token
		idFromJWT, ok := claims["adminId"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "Admin Id not found in claims",
				Error:   "Admin Id not found in claims",
			})
		}

		admin := models.NewAdminModel(h.server.IdGen, h.server.DB)
		if admin == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error initializing model",
				Error:   models.MODEL_INIT_ERROR,
			})
		}

		if _, err := admin.FindById(idFromJWT); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "Admin not found",
				Error:   err.Error(),
			})
		}

		if _, err := admin.DecryptUsername(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: encryption.DECRYPTION_ERR,
				Error:   err.Error(),
			})
		}

		accessToken, err := h.server.Auth.GenerateAdminAccessToken(authenticator.AdminOptions{
			Id:       admin.Id,
			Username: admin.Username,
			Role:     admin.Role,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: authenticator.ACCESS_TOKEN_ERR,
				Error:   err.Error(),
			})
		}

		newAccessToken = accessToken
	} else {
		// User requests refresh token
		idFromJWT, ok := claims["userId"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "User Id not found in claims",
				Error:   "User Id not found in claims",
			})
		}

		user := models.NewUserModel(h.server.IdGen, h.server.DB)
		if user == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error initializing model",
				Error:   models.MODEL_INIT_ERROR,
			})
		}

		if _, err := user.FindById(idFromJWT); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "User not found",
				Error:   err.Error(),
			})
		}

		if _, err := user.DecryptName(h.server.Encrypt.EncryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: encryption.DECRYPTION_ERR,
				Error:   err.Error(),
			})
		}
		if _, err := user.DecryptEmail(h.server.Encrypt.EncryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: encryption.DECRYPTION_ERR,
				Error:   err.Error(),
			})
		}

		accessToken, err := h.server.Auth.GenerateUserAccessToken(authenticator.UserOptions{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: authenticator.ACCESS_TOKEN_ERR,
				Error:   err.Error(),
			})
		}

		newAccessToken = accessToken
	}

	// As protection against my dumb self, check first if newAccessToken is an empty string
	if newAccessToken == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error generating new access token",
			Error:   "Error generating new access token",
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"accessToken": newAccessToken,
	})
}
