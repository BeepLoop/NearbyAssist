package handlers

import (
	"nearbyassist/internal/models"
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
	// TODO: Handle admin login
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Admin login",
	})
}

func (h *authHandler) HandleRegister(c echo.Context) error {
	user := models.NewUserModel()
	err := c.Bind(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userId, err := h.server.DB.NewUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	accessToken, err := h.server.Auth.GenerateAccessToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	session := models.NewSessionModel(userId, refreshToken)
	if _, err = h.server.DB.NewSession(session); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"userId":       userId,
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

	accessToken, err := h.server.Auth.GenerateAccessToken(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken, err := h.server.Auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	session := models.NewSessionModel(model.Id, refreshToken)
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
	refreshTokenModel := models.NewRefreshTokenModel()
	err := c.Bind(refreshTokenModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(refreshTokenModel); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session, err := h.server.DB.FindActiveSessionByToken(refreshTokenModel.Token)
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
	refreshTokenModel := models.NewRefreshTokenModel()
	err := c.Bind(refreshTokenModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(refreshTokenModel); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if blacklist, _ := h.server.DB.FindBlacklistedToken(refreshTokenModel.Token); blacklist != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Token blacklisted")
	}

	session, err := h.server.DB.FindActiveSessionByToken(refreshTokenModel.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, err := h.server.DB.FindUserById(session.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	accessToken, err := h.server.Auth.GenerateAccessToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"accessToken": accessToken,
	})
}
