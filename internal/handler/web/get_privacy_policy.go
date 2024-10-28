package web

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func GetPrivacyPolicy(c echo.Context) error {
	page := pages.PrivacyPolicy()
	return page.Render(context.Background(), c.Response().Writer)
}
