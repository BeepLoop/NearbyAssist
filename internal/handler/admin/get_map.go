package admin

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *adminHandler) GetMap(c echo.Context) error {
	markers := make([]pages.Marker, 0)

	page := pages.Map(markers)
	return page.Render(context.Background(), c.Response().Writer)
}
