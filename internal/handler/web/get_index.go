package web

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func GetIndex(c echo.Context) error {
	page := pages.Index()
	return page.Render(context.Background(), c.Response().Writer)
}
