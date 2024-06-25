package handlers

import (
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/models"
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

func (h *uploadHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Transaction base route",
	})
}

func (h *uploadHandler) HandleNewServicePhoto(c echo.Context) error {
	params, err := utils.GetUploadParams(c, "vendorId", "serviceId")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if service, _ := h.server.DB.FindServiceById(params["serviceId"]); service == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Service not found")
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)
		filename, err := handler.SavePhoto(file, h.server.Storage.SaveServicePhoto)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		uploadData := models.NewServicePhotoModel(params["vendorId"], params["serviceId"], filename)
		_, err = h.server.DB.NewServicePhoto(uploadData)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message": "Files uploaded successfully",
	})
}

func (h *uploadHandler) HandleNewProofPhoto(c echo.Context) error {
	params, err := utils.GetUploadParams(c, "applicationId", "applicantId")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if application, _ := h.server.DB.FindApplicationById(params["applicationId"]); application == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Application not found")
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)
		filename, err := handler.SavePhoto(file, h.server.Storage.SaveApplicationProof)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		uploadData := models.NewApplicationProofModel(params["applicationId"], params["applicantId"], filename)
		_, err = h.server.DB.NewApplicationProof(uploadData)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message": "Files uploaded successfully",
	})
}
