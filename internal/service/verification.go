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
	jwt       auth.Authenticator
	encryptor auth.Encryption
	fs        fs.FileStorage
}

func NewVerificationService(store verification.VerificationStore, jwt auth.Authenticator, encryptor auth.Encryption, fs fs.FileStorage) *VerificationService {
	return &VerificationService{
		store:     store,
		jwt:       jwt,
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

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	encryptedName, err := s.encryptor.EncryptString(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting name",
			Error:   auth.ENCRYPTION_ERR,
		})
	}

	encryptedAddress, err := s.encryptor.EncryptString(address)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting name",
			Error:   auth.ENCRYPTION_ERR,
		})
	}

	encryptedIdNumber, err := s.encryptor.EncryptString(idNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encrypting name",
			Error:   auth.ENCRYPTION_ERR,
		})
	}

	req := new(models.IdentityVerificationModel)
	req.UserId = userId
	req.Name = encryptedName
	req.Address = encryptedAddress
	req.IdType = idType
	req.IdNumber = encryptedIdNumber

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

	verificationId, err := s.store.Create(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating identity verification",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"verification": verificationId,
	})
}
