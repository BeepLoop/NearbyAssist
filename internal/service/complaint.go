package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/store/complaint"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ComplaintService struct {
	store     complaint.ComplaintStore
	encryptor auth.Encryption
	fs        fs.FileStorage
}

func NewComplaintService(store complaint.ComplaintStore, encryptor auth.Encryption, fs fs.FileStorage) *ComplaintService {
	return &ComplaintService{
		store:     store,
		encryptor: encryptor,
		fs:        fs,
	}
}

func (s *ComplaintService) SystemComplaint(c echo.Context) error {
	title := c.FormValue("title")
	detail := c.FormValue("detail")

	req := new(request.SystemComplaintPayload)
	req.Title = title
	req.Detail = detail

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
	}

	newComplaint := new(models.SystemComplaintModel)
	newComplaint.Title = req.Title
	newComplaint.Detail = req.Detail

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	for _, file := range files {
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error reading file",
				Error:   err.Error(),
			})
		}

		cipher, err := s.encryptor.EncryptFile(bytes)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error encrypting file",
				Error:   auth.ENCRYPTION_ERR,
			})
		}

		fileData := fs.File{
			Data:     cipher,
			Category: fs.SYS_COMPLAINT_DIR,
		}
		url, err := s.fs.SaveFile(fileData)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error saving photo",
				Error:   err.Error(),
			})
		}

		newComplaint.Images = append(newComplaint.Images, url)
	}

	if cipher, err := s.encryptor.EncryptString(req.Title); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting title",
			Error:   auth.ENCRYPTION_ERR,
		})
	} else {
		newComplaint.Title = cipher
	}

	if cipher, err := s.encryptor.EncryptString(req.Detail); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting title",
			Error:   auth.ENCRYPTION_ERR,
		})
	} else {
		newComplaint.Detail = cipher
	}

	complaintId, err := s.store.CreateSystemComplaint(newComplaint)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating system complaint",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"complaintId": complaintId,
	})
}

func (s *ComplaintService) VendorComplaint(c echo.Context) error {
	// TODO: Implement filing vendor complaint/report

	return c.JSON(http.StatusNoContent, nil)
}
