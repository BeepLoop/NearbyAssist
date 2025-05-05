package complaint

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	complaint_service "nearbyassist/internal/service/complaint"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type complaintHandler struct {
	complaintService *complaint_service.Service
}

func NewHandler(complaintService *complaint_service.Service) *complaintHandler {
	return &complaintHandler{complaintService: complaintService}
}

func (h *complaintHandler) CreateBugReport(c echo.Context) error {
	title := c.FormValue("title")
	detail := c.FormValue("detail")

	req := new(request.BugReportPayload)
	req.Title = title
	req.Detail = detail

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
	}

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	if err := h.complaintService.CreateBugReport(req, files); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating bug report",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *complaintHandler) ReportUser(c echo.Context) error {
	userId := c.FormValue("userId")
	category := c.FormValue("category")
	bookingId := c.FormValue("bookingId")
	reason := c.FormValue("reason")
	detail := c.FormValue("detail")

	req := &request.ReportUserPayload{
		UserId:    userId,
		Category:  category,
		BookingId: bookingId,
		Reason:    reason,
		Detail:    detail,
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
	}

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.complaintService.ReportUser(bearerToken, req, files); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error reporting user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
