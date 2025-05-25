package tag

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	tag_service "nearbyassist/internal/service/tag"
	"nearbyassist/internal/utils"
	"net/http"
	"slices"

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

	resp := slices.AppendSeq(
		make([]string, 0),
		utils.Map(tags, func(t *models.TagModel) string {
			return t.Title
		}),
	)

	return c.JSON(http.StatusOK, utils.Mapper{
		"tags": resp,
	})
}

func (h *tagHandler) GetExpertiseList(c echo.Context) error {
	list, err := h.tagService.GetExpertiseList()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving expertise list",
			Error:   err.Error(),
		})
	}

	resp := slices.AppendSeq(
		make([]response.Expertise, 0),
		utils.Map(list, func(e *models.ExpertiseModel) response.Expertise {
			return response.Expertise{
				Id:    e.Id,
				Title: e.Title,
			}
		}),
	)

	return c.JSON(http.StatusOK, utils.Mapper{
		"expertise": resp,
	})
}
