package middleware

import (
	"nearbyassist/internal/authenticator"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func CheckRole(jwtChecker authenticator.Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")

			err := jwtChecker.ValidateToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			claims, err := jwtChecker.GetClaims(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			role, ok := claims["role"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "Unkown user")
			}

			// Check if the user is accessing admin-only route
			url := c.Request().URL.String()
			iAdminRoute := strings.HasPrefix(url, "/admin")
			if iAdminRoute && role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "Unauthorized access")
			}

			return next(c)
		}
	}
}
