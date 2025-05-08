package complaint

import (
	"fmt"
	"nearbyassist/internal/config"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/pdfgenerator"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *complaintHandler) DownloadPDF(c echo.Context) error {
	reportId := c.Param("reportId")
	url := fmt.Sprintf("%s/admin/complaints/users/%s/pdf", config.Instance.DOMAIN, reportId)

	pdf, err := pdfgenerator.GeneratePDF(url)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "PDF Generation failed",
			Error:   err.Error(),
		})
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s.pdf"`, reportId))
	return c.Blob(http.StatusOK, "application/pdf", pdf)
}
