package settingshandler

import (
	"context"
	"fmt"
	"nearbyassist/internal/config/setting"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/activitylog"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/settings"
	"net/http"
	"time"

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

	behavior := setting.New().Values.SearchBehavior

	page := settings.Settings(*admin, flash, string(behavior))
	return page.Render(c.Request().Context(), c.Response().Writer)
}

func (h *handler) ResetSSE(c echo.Context) error {
	sse.New().SetValues(h.db)

	utils.SetFlashMessage(c, "success", "Server Sent Events (SSE) data reset")

	admin, _ := utils.GetAdminFromSession(c)

	activity := activitylog.Input{
		AdminId: admin.Id,
		Action:  activitylog.ACTION_RESETTED_SSE,
	}
	activitylog.MustGetInstance().Create(activity)

	return c.Redirect(http.StatusSeeOther, "/admin/settings")
}

func (h *handler) RemindScheduled(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            *
        FROM
            Booking
        WHERE
            scheduledAt IS NULL AND DATE(scheduledAt) = CURDATE() + INTERVAL 1 DAY
    `
	bookings := make([]*models.BookingModel, 0)
	if err := h.db.SelectContext(ctx, &bookings, query); err != nil {
		fmt.Println("error retrieving tomorrows' bookings: ", err.Error())
		utils.SetFlashMessage(c, "error", "Error occurred while retrieving bookings scheduled for tomorrow")
		return c.Redirect(http.StatusSeeOther, "/admin/settings")
	}

	if ctx.Err() == context.DeadlineExceeded {
		utils.SetFlashMessage(c, "error", "Database took too long to respond")
		return c.Redirect(http.StatusSeeOther, "/admin/settings")
	}

	sent := 0

	oneSignal := notification_service.MustGetInstance()
	for _, booking := range bookings {
		heading := "Upcoming Booking Tomorrow"
		content := "You have a booking scheduled for tomorrow. Get ready to deliver your service on time."

		if err := oneSignal.NewUrgentNotification(booking.VendorId, heading, content); err != nil {
			fmt.Println("error notifying vendor %s, error: ", booking.VendorId, err.Error())
			continue
		}

		sent++
	}

	if err := utils.SetFlashMessage(c, "success", fmt.Sprintf("Sent reminder to %d vendors", sent)); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/settings?error=sent_remined_to_vendors")
	}

	admin, _ := utils.GetAdminFromSession(c)

	activity := activitylog.Input{
		AdminId: admin.Id,
		Action:  activitylog.ACTION_NOTIFIED_VENDORS,
	}
	activitylog.MustGetInstance().Create(activity)

	return c.Redirect(http.StatusSeeOther, "/admin/settings")
}

func (h *handler) UpdateSeachBehavior(c echo.Context) error {
	behavior := c.FormValue("behavior")

	if !setting.IsValidSearchBehavior(behavior) {
		utils.SetFlashMessage(c, "error", "Invalid search behavior provided")
		return c.Redirect(http.StatusSeeOther, "/admin/settings")
	}

	setting.New().Values.SearchBehavior = setting.SearchBehavior(behavior)

	if err := utils.SetFlashMessage(c, "success", "search behavior updated"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/settings?success=search_behavior_changed")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/settings")
}
