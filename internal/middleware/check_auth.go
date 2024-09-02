package middleware

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func CheckAuth(jwtChecker authenticator.Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if err := jwtChecker.ValidateToken(token); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
					Message: "Invalid token",
					Error:   err.Error(),
				})
			}

			return next(c)
		}
	}
}
