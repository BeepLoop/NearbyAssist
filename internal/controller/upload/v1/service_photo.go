package upload

import (
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

	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	files := form.File["files"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no files attached")
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
