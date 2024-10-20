package tag

import (
	"nearbyassist/internal/models"
	tag_service "nearbyassist/internal/service/tag"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type tagHandler struct {
	tagService *tag_service.Service
}

func NewHandler(tagService *tag_service.Service) *tagHandler {
	return &tagHandler{tagService: tagService}
}

func (h *tagHandler) GetTags(c echo.Context) error {
	tags, err := h.tagService.GetTags()
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
