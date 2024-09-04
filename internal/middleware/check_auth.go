package middleware

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckAuth(jwtChecker authenticator.Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
            token := c.Request().Header.Get("Authorization")[len("Bearer "):]
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
