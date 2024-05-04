package handlers

import (
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
	adminModel := models.NewAdminModel()
	err := c.Bind(adminModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(adminModel); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO: Handle validating credentials
	admin, err := h.server.DB.FindAdminByUsername(adminModel.Username)
	if admin == nil || err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Account not found")
	}

	// TODO: Implement better password validation with encryption
	if adminModel.Password != admin.Password {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid credentials")
	}

	accessToken, err := h.server.Auth.GenerateAdminAccessToken(admin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: Implement session management for admin
	session := models.NewSessionModel(refreshToken)
	if _, err := h.server.DB.NewSession(session); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"adminId":      admin.Id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *authHandler) HandleLogin(c echo.Context) error {
	model := models.NewUserModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(model); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.server.DB.FindUserByEmail(model.Email)
	if err != nil {
		userId, err := h.server.DB.NewUser(model)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		model.Id = userId
	} else {
		model.Id = user.Id
	}

	accessToken, err := h.server.Auth.GenerateUserAccessToken(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	session := models.NewSessionModel(refreshToken)
	if _, err := h.server.DB.NewSession(session); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"userId":       model.Id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *authHandler) HandleLogout(c echo.Context) error {
	req := &request.Logout{}
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session, err := h.server.DB.FindActiveSessionByToken(req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Session not found")
	}

	if err := h.server.DB.LogoutSession(session.Id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := h.server.DB.BlacklistToken(session.Token); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Logout successful",
	})
}

func (h *authHandler) HandleTokenRefresh(c echo.Context) error {
	req := &request.RefreshToken{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if blacklist, _ := h.server.DB.FindBlacklistedToken(req.Token); blacklist != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Token blacklisted")
	}

	authHeader := c.Request().Header.Get("Authorization")
	userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, err := h.server.DB.FindUserById(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	accessToken, err := h.server.Auth.GenerateUserAccessToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"accessToken": accessToken,
	})
}
