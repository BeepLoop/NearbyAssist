package handlers

import (
	"nearbyassist/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type categoryHandler struct {
	server *server.Server
}

func NewCategoryHandler(server *server.Server) *categoryHandler {
	return &categoryHandler{
		server: server,
	}
}

func (h *categoryHandler) HandleCategories(c echo.Context) error {
	categories, err := h.server.DB.FindAllCategory()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}
