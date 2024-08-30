package handlers

import (
	"nearbyassist/internal/encryption"
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type complaintHandler struct {
	server *server.Server
}

func NewComplaintHandler(server *server.Server) *complaintHandler {
	return &complaintHandler{
		server: server,
	}
}

func (h *complaintHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Complaints base route",
	})
}

func (h *complaintHandler) HandleSystemComplaint(c echo.Context) error {
	title := c.FormValue("title")
	detail := c.FormValue("detail")

	complaint := models.NewSystemComplaintModel(h.server.IdGen, h.server.DB)
	if complaint == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	} else {
		complaint.Title = title
		complaint.Detail = detail
	}

	if _, err := complaint.EncryptTitle(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	if _, err := complaint.EncryptDetail(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := complaint.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)
		url, err := handler.SavePhoto(file, h.server.Storage.SaveSystemComplaint)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		imageData := models.SystemComplaintData{ComplaintId: complaint.Id, Url: url}
		image := models.NewSystemComplaintImageWithData(imageData, h.server.IdGen, h.server.DB)
		if image == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
		}

		if _, err := image.Create(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while saving image")
		}
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message":     "System complaint created successfully",
		"complaintId": complaint.Id,
	})
}

func (h *complaintHandler) HandleSystemComplaintCount(c echo.Context) error {
	systemComplaint := models.NewSystemComplaintModel(h.server.IdGen, h.server.DB)
	if systemComplaint == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	count, err := systemComplaint.Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *complaintHandler) HandleGetSystemComplaint(c echo.Context) error {
	systemComplaint := models.NewSystemComplaintModel(h.server.IdGen, h.server.DB)
	if systemComplaint == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	complaints, err := systemComplaint.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, complaint := range complaints {
		if _, err := complaint.DecryptTitle(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
		}
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaints": complaints,
	})
}

func (h *complaintHandler) HandleGetSystemComplaintById(c echo.Context) error {
	complaintId := c.Param("complaintId")
	if complaintId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Complaint ID must be a number")
	}

	systemComplaint := models.NewSystemComplaintModel(h.server.IdGen, h.server.DB)
	if systemComplaint == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := systemComplaint.FindById(complaintId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Complaint not found")
	}

	if _, err := systemComplaint.DecryptTitle(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
	}

	if _, err := systemComplaint.DecryptDetail(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
	}

	images, err := systemComplaint.GetPhotos()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaint": systemComplaint,
		"images":    images,
	})
}

func (h *complaintHandler) HandleVendorComplaint(c echo.Context) error {
	// TODO: implement filing a complaint for a vendor
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Vendor complaint route",
	})
}
