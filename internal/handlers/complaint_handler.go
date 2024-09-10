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
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	} else {
		complaint.Title = title
		complaint.Detail = detail
	}

	if _, err := complaint.EncryptTitle(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting title",
			Error:   encryption.ENCRYPTION_ERR,
		})
	}

	if _, err := complaint.EncryptDetail(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting detail",
			Error:   encryption.ENCRYPTION_ERR,
		})
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	if _, err := complaint.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating system complaint",
			Error:   err.Error(),
		})
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)
		url, err := handler.SavePhoto(file, h.server.Storage.SaveSystemComplaint)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error saving photo",
				Error:   err.Error(),
			})
		}

		imageData := models.SystemComplaintData{ComplaintId: complaint.Id, Url: url}
		image := models.NewSystemComplaintImageWithData(imageData, h.server.IdGen, h.server.DB)
		if image == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error initializing image model",
				Error:   models.MODEL_INIT_ERROR,
			})
		}

		if _, err := image.Create(); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error creating image",
				Error:   err.Error(),
			})
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
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	count, err := systemComplaint.Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting system complaint count",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *complaintHandler) HandleGetSystemComplaint(c echo.Context) error {
	systemComplaint := models.NewSystemComplaintModel(h.server.IdGen, h.server.DB)
	if systemComplaint == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	params := utils.ParseQuery(c.QueryString())
	complaints, err := systemComplaint.FindAll(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting system complaints",
			Error:   err.Error(),
		})
	}

	for _, complaint := range complaints {
		if _, err := complaint.DecryptTitle(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error decrypting title",
				Error:   encryption.DECRYPTION_ERR,
			})
		}
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaints": complaints,
	})
}

func (h *complaintHandler) HandleGetSystemComplaintById(c echo.Context) error {
	complaintId := c.Param("complaintId")
	if complaintId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Complaint ID is required",
			Error:   "Complaint ID is required",
		})
	}

	systemComplaint := models.NewSystemComplaintModel(h.server.IdGen, h.server.DB)
	if systemComplaint == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := systemComplaint.FindById(complaintId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Complaint not found",
			Error:   err.Error(),
		})
	}

	if _, err := systemComplaint.DecryptTitle(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error decrypting title",
			Error:   encryption.DECRYPTION_ERR,
		})
	}

	if _, err := systemComplaint.DecryptDetail(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error decrypting detail",
			Error:   encryption.DECRYPTION_ERR,
		})
	}

	images, err := systemComplaint.GetPhotos()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting complaint images",
			Error:   err.Error(),
		})
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
