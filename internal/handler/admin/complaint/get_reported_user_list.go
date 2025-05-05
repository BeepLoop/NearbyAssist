package complaint

import (
	"context"
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/complaints"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) GetReportList(c echo.Context) error {
	params := c.QueryParams()
	flash, _, _ := utils.RetrieveFlashMessage(c)
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	reports := make([]dto.ReportItem, 0)
	if params.Has("query") {
		// Perform a search
	} else {
		limit, _ := strconv.Atoi(params.Get("limit"))
		if limit == 0 {
			limit = DEFAULT_LIMIT
		}

		offset, _ := strconv.Atoi(params.Get("offset"))
		if offset == 0 {
			offset = DEFAULT_OFFSET
		}

		data, err := h.complaintService.GetReportList(limit, offset)
		if err != nil {
			fmt.Println(err.Error())
			page := pages.ReportedUserList(*admin, make([]dto.ReportItem, 0), flash)
			return page.Render(context.Background(), c.Response().Writer)
		}
		reports = data
	}

	page := pages.ReportedUserList(*admin, reports, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
