package handlers

import (
	"nearbyassist/internal/models"
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
	tag := models.NewTagModel(h.server.IdGen, h.server.DB)
	if tag == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error initializing model",
            Error:   models.MODEL_INIT_ERROR,
        })
	}

	tags, err := tag.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error getting tags",
            Error:   err.Error(),
        })
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"tags": tags,
	})
}
