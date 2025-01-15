package web

import (
	"context"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func GetAccountDeletionInstructions(c echo.Context) error {
	page := pages.AccountDeletionInstructions()
	return page.Render(context.Background(), c.Response().Writer)
}
