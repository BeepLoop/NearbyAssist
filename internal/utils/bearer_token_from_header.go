package utils

import "github.com/labstack/echo/v4"

func BearerTokenFromHeader(c echo.Context) string {
	return c.Request().Header.Get("Authorization")[len("Bearer "):]
}
