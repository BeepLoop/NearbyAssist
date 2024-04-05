package upload

import (
	service_query "nearbyassist/internal/db/query/service"
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ServicePhoto(c echo.Context) error {
	params, err := utils.GetUploadParams(c, "vendorId", "serviceId")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	doesVendorExists := service_query.DoesServiceExists(params["serviceId"], params["vendorId"])
	if !doesVendorExists {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor not found")
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var uploadId int
	for _, file := range files {
		handler := filehandler.NewServicePhoto(params["vendorId"], params["serviceId"], file)
		uploadId, err = handler.Upload()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"uploadId": uploadId,
	})
}
