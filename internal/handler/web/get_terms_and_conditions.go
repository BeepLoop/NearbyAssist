package web

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func GetTermsAndConditions(c echo.Context) error {
	page := pages.TermsAndConditions()
	return page.Render(context.Background(), c.Response().Writer)
}
