package expertise

import (
	"context"
	"nearbyassist/internal/models"
	expertise_service "nearbyassist/internal/service/expertise"
	pages "nearbyassist/views/pages/expertise"

	"github.com/labstack/echo/v4"
)

type expertiseHandler struct {
	expertService *expertise_service.Service
}

func NewHandler(expertService *expertise_service.Service) *expertiseHandler {
	return &expertiseHandler{
		expertService: expertService,
	}
}

func (h *expertiseHandler) GetIndex(c echo.Context) error {
	results, err := h.expertService.GetAllExpertise()
	if err != nil {
		page := pages.Expertise(make([]models.ExpertiseModel, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.ExpertiseModel, 0)
	for _, entry := range results {
		data = append(data, models.ExpertiseModel{
			Model: entry.Model,
			Title: entry.Title,
			Tags:  entry.Tags,
		})
	}

	page := pages.Expertise(data)
	return page.Render(context.Background(), c.Response().Writer)
}
