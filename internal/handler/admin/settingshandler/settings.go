package settingshandler

import (
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/settings"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *handler {
	return &handler{db: db}
}

func (h *handler) GetSettingsPage(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	page := settings.Settings(*admin, flash)
	return page.Render(c.Request().Context(), c.Response().Writer)
}

func (h *handler) ResetSSE(c echo.Context) error {
	sse.New().SetValues(h.db)

	utils.SetFlashMessage(c, "success", "Server Sent Events (SSE) data reset")

	return c.Redirect(http.StatusSeeOther, "/admin/settings")
}
