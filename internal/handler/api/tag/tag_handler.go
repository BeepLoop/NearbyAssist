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

func (h *tagHandler) GetExpertise(c echo.Context) error {
	expertiseWithTags, err := h.tagService.GetExpertise()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting expertise",
			Error:   err.Error(),
		})
	}

	response := []struct {
		Id    string `json:"id"`
		Title string `json:"title"`
		Tags  []struct {
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"tags"`
	}{}

	for _, entry := range expertiseWithTags {
		tags := make([]struct {
			Id    string `json:"id"`
			Title string `json:"title"`
		}, 0)
		for _, tag := range entry.Tags {
			tags = append(tags, struct {
				Id    string `json:"id"`
				Title string `json:"title"`
			}{
				Id:    tag.Id,
				Title: tag.Title,
			})
		}

		response = append(response, struct {
			Id    string `json:"id"`
			Title string `json:"title"`
			Tags  []struct {
				Id    string `json:"id"`
				Title string `json:"title"`
			} `json:"tags"`
		}{
			Id:    entry.Id,
			Title: entry.Title,
			Tags:  tags,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"expertises": response,
	})
}
