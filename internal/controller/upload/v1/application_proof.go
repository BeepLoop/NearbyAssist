package upload

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func VendorApplicationProof(c echo.Context) error {
	params, err := utils.GetUploadParams(c, "applicationId", "applicantId")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	doesApplicationExists := vendor_query.DoesApplicationExists(params["applicationId"], params["applicantId"])
	if !doesApplicationExists {
		return echo.NewHTTPError(http.StatusBadRequest, "application not found")
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var uploadId int
	for _, file := range files {
		handler := filehandler.NewApplicationProof(params["applicationId"], params["applicantId"], file)
		uploadId, err = handler.Upload()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"uploadId": uploadId,
	})
}
