package map_handler

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/map"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *mapHandler) GetMap(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	services := make([]models.ServiceModel, 0)
	tags := make([]string, 0)

	pageData := models.MapPageDataModel{
		Services: services,
		Tags:     tags,
	}

	params := c.QueryParams()
	if params.Has("query") {
		services, err := h.mapService.GetServices(params.Get("query"))
		if err != nil {
			page := pages.Map(*admin, pageData, flash)
			return page.Render(context.Background(), c.Response().Writer)
		}

		for _, service := range services {
			pageData.Services = append(pageData.Services, models.ServiceModel{
				Model:       service.Model,
				Address:     service.Address,
				VendorId:    service.VendorId,
				Title:       service.Title,
				Description: service.Description,
				Rate:        service.Rate,
			})
		}
	}

	result, err := h.tagService.GetTags()
	if err != nil {
		page := pages.Map(*admin, pageData, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	for _, tag := range result {
		pageData.Tags = append(pageData.Tags, tag.Title)
	}

	page := pages.Map(*admin, pageData, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
