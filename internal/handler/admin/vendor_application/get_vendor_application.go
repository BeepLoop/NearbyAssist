package application

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *applicationHandler) GetVendorApplication(c echo.Context) error {
	applications, err := h.applicationService.GetApplications()
	if err != nil {
		page := pages.Applications(make([]models.ApplicationModel, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.Applications(applications)
	return page.Render(context.Background(), c.Response().Writer)
}
