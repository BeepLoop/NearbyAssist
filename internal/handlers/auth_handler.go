package handlers

import (
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"

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

	err = c.Validate(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userId, err := h.server.DB.NewUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: Generate JWT Token, not really sure
	token, err := utils.GenerateJwt(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: Implement better session handling
	session := models.NewSessionModel(userId, token)
	sessionId, err := h.server.DB.NewSession(session)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"userId":    userId,
		"sessionId": sessionId,
	})
}

func (h *authHandler) HandleLogin(c echo.Context) error {
	model := models.NewUserModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.server.DB.FindUserByEmail(model.Email)
	if err != nil {
		fmt.Println("user  not found, trying to register")
		userId, err := h.server.DB.NewUser(model)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		model.Id = userId
	} else {
		model.Id = user.Id
	}

	// TODO: Generate JWT Token, not really sure
	token, err := utils.GenerateJwt(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: Implement better session handling
	session := models.NewSessionModel(model.Id, token)
	sessionId, err := h.server.DB.NewSession(session)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"userId":    model.Id,
		"sessionId": sessionId,
	})
}

func (h *authHandler) HandleLogout(c echo.Context) error {
	// TODO: Handle logout
	cookie, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	sessionId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if session, _ := h.server.DB.FindSessionById(sessionId); session == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Session not found")
	}

	if err := h.server.DB.LogoutSession(sessionId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Logout",
	})
}
