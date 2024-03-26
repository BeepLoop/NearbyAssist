package service

import (
	service_query "nearbyassist/internal/db/query/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetOwnerServices(c echo.Context) error {
	ownerId := c.Param("ownerId")
	if ownerId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "owner ID missing",
		})
	}
	id, err := strconv.Atoi(strings.ReplaceAll(ownerId, "/", ""))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "owner ID must be a number",
		})
	}

	services, err := service_query.GetOwnerServices(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, services)
}
