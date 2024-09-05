package middleware

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func CheckRole(jwtChecker authenticator.Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")[len("Bearer "):]
			if err := jwtChecker.ValidateToken(token); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
					Message: "Invalid token",
					Error:   err.Error(),
				})
			}

			claims, err := jwtChecker.GetClaims(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
					Message: "Error getting claims",
					Error:   err.Error(),
				})
			}

			role, ok := claims["role"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, models.Error{
					Message: "Role not found",
					Error:   "Role not found",
				})
			}

			// Check if the user is accessing admin-only route
			url := c.Request().URL.String()
			iAdminRoute := strings.Contains(url, "/admin")
			if iAdminRoute && role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, models.Error{
					Message: "Unauthorized access",
					Error:   "Unauthorized access",
				})
			}

			return next(c)
		}
	}
}
