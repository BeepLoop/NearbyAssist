package middleware

import (
	"nearbyassist/internal/authenticator"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func CheckAuth(jwtChecker authenticator.Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")

			err := jwtChecker.ValidateToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
			}

			return next(c)
		}
	}
}
