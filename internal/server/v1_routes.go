package server

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HandleVersionOneRoutes(r *echo.Group) {

	r.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	r.POST("/register", s.HandleRegister)

	r.POST("/login", s.HandleLogin)

	r.GET("/locations", func(c echo.Context) error {
		locations, err := db.GetLocations()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"locations": locations,
		})
	})

	r.POST("/locations/register", func(c echo.Context) error {
		var location types.Location
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
		var position types.Position
		err := c.Bind(&position)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid request",
			})
		}

		locations, err := db.SearchVicinity(position)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, locations)
	})
}
