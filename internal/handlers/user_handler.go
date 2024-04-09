package handlers

import (
	"nearbyassist/internal/db/models"
	"nearbyassist/internal/server"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	server *server.Server
}

func NewUserHandler(server *server.Server) *userHandler {
	return &userHandler{
		server: server,
	}
}

func (h *userHandler) HandleGetUser(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	userModel := models.NewUserModel()
	user, err := userModel.FindById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) HandleCount(c echo.Context) error {
	model := models.NewUserModel()

	count, err := model.Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"userCount": count,
	})
}
