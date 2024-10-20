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

func (h *complaintHandler) CreateSystemComplaint(c echo.Context) error {
	title := c.FormValue("title")
	detail := c.FormValue("detail")

	req := new(request.SystemComplaintPayload)
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

	complaintId, err := h.complaintService.CreateSystemComplaint(req, files)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating system complaint",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"complaintId": complaintId,
	})
}

func (h *complaintHandler) ReportVendor(c echo.Context) error {
	// TODO: Implement filing vendor complaint/report
	return c.JSON(http.StatusNoContent, nil)
}
