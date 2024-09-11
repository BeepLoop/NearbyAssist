package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/store/tag"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TagService struct {
	store tag.TagStore
}

func NewTagService(store tag.TagStore) *TagService {
	return &TagService{
		store: store,
	}
}

func (s *TagService) BaseRoute(c echo.Context) error {
	tags, err := s.store.FindAll()
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
