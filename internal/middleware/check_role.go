package middleware

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "Missing token",
				Error:   "Missing token",
			})
		}

		token = token[len("Bearer "):]
		if err := auth.ValidateToken(token); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "Invalid token",
				Error:   err.Error(),
			})
		}

		claims, err := auth.GetClaims(token)
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
		_ = role

		// Check if the user is accessing admin-only route
		// url := c.Request().URL.String()
		// iAdminRoute := strings.Contains(url, "/admin")
		// if iAdminRoute && role != "admin" {
		// 	return echo.NewHTTPError(http.StatusForbidden, models.Error{
		// 		Message: "Unauthorized access",
		// 		Error:   "Unauthorized access",
		// 	})
		// }

		return next(c)
	}
}
