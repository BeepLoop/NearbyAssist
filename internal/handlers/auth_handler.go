package handlers

import (
	"fmt"
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

	err = c.Validate(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userId, err := h.server.DB.NewUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: Generate JWT Token, not really sure

	// TODO: Create new session

	return c.JSON(http.StatusCreated, utils.Mapper{
		"userId": userId,
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

	// TODO: Create new session

	return c.JSON(http.StatusCreated, utils.Mapper{
		"userId": model.Id,
	})
}

func (h *authHandler) HandleLogout(c echo.Context) error {
	// TODO: Handle logout
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Logout",
	})
}
