package health

import (
	health_service "nearbyassist/internal/service/health"
	"net/http"

	"github.com/labstack/echo/v4"
)

type healthHandler struct {
	healthService *health_service.Service
}

func NewHandler(healthService *health_service.Service) *healthHandler {
	return &healthHandler{healthService: healthService}
}

func (h *healthHandler) CheckHealth(c echo.Context) error {
	health := h.healthService.CheckHealth()

	return c.JSON(http.StatusOK, health)
}
