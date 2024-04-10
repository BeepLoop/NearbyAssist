package handlers

import (
	"nearbyassist/internal/db/models"
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type uploadHandler struct {
	server *server.Server
}

func NewUploadHandler(server *server.Server) *uploadHandler {
	return &uploadHandler{
		server: server,
	}
}

func (h *uploadHandler) HandleNewServicePhoto(c echo.Context) error {
	params, err := utils.GetUploadParams(c, "vendorId", "serviceId")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	model := models.NewServiceModel(h.server.DB)
	if service, _ := model.FindById(params["serviceId"]); service == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Service not found")
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, file := range files {
		model := models.NewServicePhotoModel(params["vendorId"], params["serviceId"], h.server.DB)
		handler := filehandler.NewFileHandler(model)

		_, err := handler.SaveFile(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Files uploaded successfully",
	})
}

func (h *uploadHandler) HandleNewProofPhoto(c echo.Context) error {
	params, err := utils.GetUploadParams(c, "applicationId", "applicantId")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	model := models.NewApplicationModel(h.server.DB)
	if application, _ := model.FindById(params["applicationId"]); application == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Application not found")
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, file := range files {
		model := models.NewApplicationProofModel(params["applicationId"], params["applicantId"], h.server.DB)
		handler := filehandler.NewFileHandler(model)

		_, err := handler.SaveFile(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Files uploaded successfully",
	})
}
