package map_handler

import (
	"context"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/map"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *mapHandler) GetMap(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	query := utils.Ternary(c.QueryParams().Has("query"), c.QueryParams().Get("query"), "")
	mapdata, _ := h.mapService.GetMapData(query)

	page := pages.Map(*admin, *mapdata, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
