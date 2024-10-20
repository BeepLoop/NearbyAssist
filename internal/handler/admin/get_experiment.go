package admin

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *adminHandler) GetExperiment(c echo.Context) error {
	page := pages.Experiment()
	return page.Render(context.Background(), c.Response().Writer)
}
