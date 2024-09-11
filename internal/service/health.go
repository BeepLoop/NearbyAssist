package handler

import (
	"nearbyassist/internal/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) BaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"health": "ok",
		"time":   time.Now(),
	})
}
