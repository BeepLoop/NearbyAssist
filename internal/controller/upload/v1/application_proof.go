package upload

import (
	"nearbyassist/internal/db/models"
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

	model := models.NewApplicationModel()
	if application, _ := model.FindById(params["applicationId"]); application == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Application not found")
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, file := range files {
		model := models.NewApplicationProofModel(params["applicationId"], params["applicantId"])
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
