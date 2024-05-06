package handlers

import (
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type tagHandler struct {
	server *server.Server
}

func NewTagHandler(server *server.Server) *tagHandler {
	return &tagHandler{
		server: server,
	}
}

func (h *tagHandler) HandleGetTags(c echo.Context) error {
	tags, err := h.server.DB.FindAllTags()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"tags": tags,
	})
}
