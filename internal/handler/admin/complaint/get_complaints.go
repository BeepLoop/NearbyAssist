package complaint

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) GetComplaints(c echo.Context) error {
	complaints, err := h.complaintService.GetComplaints()
	if err != nil {
		page := pages.Complaints(make([]models.ComplaintModel, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.Complaints(complaints)
	return page.Render(context.Background(), c.Response().Writer)
}
