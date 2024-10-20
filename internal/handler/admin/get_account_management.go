package admin

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *adminHandler) GetAccountManagement(c echo.Context) error {
	page := pages.AccountManagement()
	return page.Render(context.Background(), c.Response().Writer)
}
