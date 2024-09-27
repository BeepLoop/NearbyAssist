package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/store/verification"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type VerificationService struct {
	store     verification.VerificationStore
	encryptor auth.Encryption
	fs        fs.FileStorage
}

func NewVerificationService(store verification.VerificationStore, encryptor auth.Encryption, fs fs.FileStorage) *VerificationService {
	return &VerificationService{
		store:     store,
		encryptor: encryptor,
		fs:        fs,
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

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	for _, file := range files {
		// Read bytes
		bytes, err := utils.FileToBytes(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error reading file",
				Error:   err.Error(),
			})
		}

		// Encrypt the file
		cipher, err := s.encryptor.EncryptFile(bytes)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error encrypting file",
				Error:   err.Error(),
			})
		}

		switch file.Filename {
		case "frontId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_FRONT,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving front ID",
					Error:   err.Error(),
				})
			} else {
				req.FrontIdImageUrl = url
			}

		case "backId":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.ID_BACK,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving back ID",
					Error:   err.Error(),
				})
			} else {
				req.BackIdImageUrl = url
			}

		case "face":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.FACE,
			}
			if url, err := s.fs.SaveFile(fileData); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving face",
					Error:   err.Error(),
				})
			} else {
				req.FaceImageUrl = url
			}

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
