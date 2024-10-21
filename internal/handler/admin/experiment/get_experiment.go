package experiment

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *experimentHandler) GetExperiment(c echo.Context) error {
	page := pages.Experiment()
	return page.Render(context.Background(), c.Response().Writer)
}
