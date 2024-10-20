package web

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func GetNotFound(c echo.Context) error {
	page := pages.NotFound()
	return page.Render(context.Background(), c.Response().Writer)
}
