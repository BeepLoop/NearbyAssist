package map_handler

import (
	"context"
	"nearbyassist/views/pages/map"

	"github.com/labstack/echo/v4"
)

func (h *mapHandler) GetMap(c echo.Context) error {
	markers := make([]pages.Marker, 0)

	page := pages.Map(markers)
	return page.Render(context.Background(), c.Response().Writer)
}
