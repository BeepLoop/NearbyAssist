package upload

import (
	"nearbyassist/internal/db/models"
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

	model := models.NewServiceModel()
	if service, _ := model.FindById(params["serviceId"]); service == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Service not found")
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, file := range files {
		model := models.NewServicePhotoModel(params["vendorId"], params["serviceId"])
		handler := filehandler.NewFileHandler(model)

		_, err := handler.SaveFile(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "File uploaded successfully",
	})
}
