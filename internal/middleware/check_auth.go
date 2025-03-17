package middleware

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/core"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckAuth(jwt core.Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
					Message: "Missing token",
					Error:   "Missing token",
				})
			}

			token = token[len("Bearer "):]
			if err := jwt.ValidateToken(token); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
					Message: "Invalid token",
					Error:   err.Error(),
				})
			}

			return next(c)
		}
	}
}
