package expertise

import (
	"context"
	"nearbyassist/internal/models"
	expertise_service "nearbyassist/internal/service/expertise"
	pages "nearbyassist/views/pages/expertise"
	"net/http"

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

func (h *expertiseHandler) GetAllExpertise(c echo.Context) error {
	params := c.QueryParams()

	results := make([]*models.ExpertiseModel, 0)

	if params.Has("query") && params.Get("query") != "" {
		title := params.Get("query")

		expertise, err := h.expertService.FindExpertise(title)
		if err != nil {
			page := pages.Expertise(make([]models.ExpertiseModel, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = append(results, expertise)
	} else {
		experitises, err := h.expertService.GetAllExpertise()
		if err != nil {
			page := pages.Expertise(make([]models.ExpertiseModel, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = experitises
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

func (h *expertiseHandler) CreateExpertise(c echo.Context) error {
	title := c.FormValue("title")

	data := &models.ExpertiseModel{
		Title: title,
	}

	if _, err := h.expertService.CreateExpertise(data); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/expertise?error=creation_error")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/expertise")
}

func (h *expertiseHandler) AddTagToExpertise(c echo.Context) error {
	expertiseId := c.FormValue("expertiseId")
	title := c.FormValue("title")

	data := &models.TagModel{
		Title: title,
	}

	if _, err := h.expertService.AddTagToExpertise(expertiseId, data); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/expertise?error=creation_error")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/expertise")
}
