package web

import (
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func GetIndex(c echo.Context) error {
	page := pages.Index()
	return page.Render(c.Request().Context(), c.Response().Writer)
}
