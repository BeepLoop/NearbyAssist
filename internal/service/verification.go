package handler

import (
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/store/verification"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type VerificationService struct {
	store     verification.VerificationStore
	disk      storage.Storage
	encryptor auth.Encryption
}

func NewVerificationService(store verification.VerificationStore, encryptor auth.Encryption, disk storage.Storage) *VerificationService {
	return &VerificationService{
		store:     store,
		disk:      disk,
		encryptor: encryptor,
	}
}

func (s *VerificationService) Create(c echo.Context) error {
	name := c.FormValue("name")
	address := c.FormValue("address")
	idType := c.FormValue("idType")
	idNumber := c.FormValue("idNumber")
	if name == "" || address == "" || idType == "" || idNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Missing required fields",
			Error:   "Missing required fields",
		})
	}

	req := new(models.IdentityVerificationModel)
	req.Name = name
	req.Address = address
	req.IdType = idType
	req.IdNumber = idNumber

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(s.encryptor)

		switch file.Filename {
		case "frontId":
			url, err := handler.SavePhoto(file, s.disk.SaveFrontId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving front ID",
					Error:   err.Error(),
				})
			}
			req.FrontIdImageUrl = url

		case "backId":
			url, err := handler.SavePhoto(file, s.disk.SaveBackId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving back ID",
					Error:   err.Error(),
				})
			}
			req.BackIdImageUrl = url

		case "face":
			url, err := handler.SavePhoto(file, s.disk.SaveFace)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving face",
					Error:   err.Error(),
				})
			}
			req.FaceImageUrl = url

		default:
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "Filename parameter not supported or known",
				Error:   "Unsupported file type",
			})
		}
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	if cipher, err := s.encryptor.EncryptString(req.Name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting name",
			Error:   auth.ENCRYPTION_ERR,
		})
	} else {
		req.Name = cipher
	}

	if cipher, err := s.encryptor.EncryptString(req.Address); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting name",
			Error:   auth.ENCRYPTION_ERR,
		})
	} else {
		req.Address = cipher
	}

	if cipher, err := s.encryptor.EncryptString(req.IdNumber); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting name",
			Error:   auth.ENCRYPTION_ERR,
		})
	} else {
		req.IdNumber = cipher
	}

	if _, err := s.store.Create(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating identity verification",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"verification": req.Id,
	})
}
