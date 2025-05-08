package log_handler

import (
	"context"
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
