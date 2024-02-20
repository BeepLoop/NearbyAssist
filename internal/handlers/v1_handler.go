package handlers

import (
	"nearbyassist/internal/controller/auth/v1"
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RouteHandlerV1(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	authGroup := r.Group("/auth")
	auth.AuthHandler(authGroup)

	r.GET("/locations", func(c echo.Context) error {
		locations, err := db.GetLocations()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, locations)
	})

	r.POST("/locations/register", func(c echo.Context) error {
		var location types.LocationRegister
		err := c.Bind(&location)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid request",
			})
		}

		err = c.Validate(location)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		err = db.RegisterLocation(location)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	r.GET("/locations/vicinity", func(c echo.Context) error {
		latitude := c.QueryParam("lat")
		longitude := c.QueryParam("long")
		position := types.Position{Latitude: latitude, Longitude: longitude}

		locations, err := db.SearchVicinity(position)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, locations)
	})
}
