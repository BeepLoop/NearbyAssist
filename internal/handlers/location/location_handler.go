package location

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LocationHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	r.GET("/", func(c echo.Context) error {
		locations, err := db.GetLocations()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, locations)
	})

	r.POST("/register", func(c echo.Context) error {
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

	r.GET("/vicinity", func(c echo.Context) error {
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
