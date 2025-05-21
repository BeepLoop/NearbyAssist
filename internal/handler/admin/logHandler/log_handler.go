package log_handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/activitylog"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/logspage"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	DEFAULT_LIMIT  = 20
	DEFAULT_OFFSET = 0
)

type logHandler struct {
}

func NewHandler() *logHandler {
	return &logHandler{}
}

func (h *logHandler) GetLogs(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	params := c.QueryParams()
	limit, _ := strconv.Atoi(params.Get("limit"))
	if limit == 0 {
		limit = DEFAULT_LIMIT
	}

	offset, _ := strconv.Atoi(params.Get("offset"))
	if offset == 0 {
		offset = DEFAULT_OFFSET
	}

	logs, err := activitylog.MustGetInstance().GetAll(limit, offset, params.Get("range"))
	if err != nil {
		fmt.Println(err.Error())
		page := logspage.Logs(*admin, make([]dto.ActivityLog, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := slices.AppendSeq(
		make([]dto.ActivityLog, 0),
		utils.Map(logs, func(log *models.ActivityLogModel) dto.ActivityLog {
			return dto.ActivityLog{
				ID:            log.Id,
				AdminUsername: log.AdminUsername,
				Action:        log.Action,
				TargetEmail:   log.Target,
				CreatedAt:     utils.FormatDate(log.CreatedAt),
			}
		}),
	)

	page := logspage.Logs(*admin, data)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *logHandler) DownloadCSV(c echo.Context) error {
	params := c.QueryParams()
	logs, err := activitylog.MustGetInstance().GetAllNoLimit(params.Get("range"))
	if err != nil {
		if err := utils.SetFlashMessage(c, "error", "Failed to generate CSV File"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/logs?error=csv_file_generation_failed")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/logs")
	}

	data := [][]string{{"ID", "Admin Username", "Action", "Target", "Created At"}}
	for _, log := range logs {
		row := []string{
			log.Id,
			log.AdminUsername,
			log.Action,
			log.Target,
			utils.FormatDate(log.CreatedAt),
		}

		data = append(data, row)
	}

	ts := utils.FilenameFriendlyTimeStamp()
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s_logs.csv"`, ts))
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")

	writer := csv.NewWriter(c.Response())

	if err := writer.WriteAll(data); err != nil {
		if err := utils.SetFlashMessage(c, "error", "Failed to generate CSV File"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/logs?error=csv_file_generation_failed")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/logs")
	}

	writer.Flush()
	return nil
}
