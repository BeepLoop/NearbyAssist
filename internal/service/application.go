package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"nearbyassist/internal/store/application"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ApplicationService struct {
	store     application.ApplicationStore
	jwt       auth.Authenticator
	encryptor auth.Encryption
	fs        fs.FileStorage
}

func NewApplicationService(store application.ApplicationStore, jwt auth.Authenticator, encryptor auth.Encryption, fs fs.FileStorage) *ApplicationService {
	return &ApplicationService{
		store:     store,
		jwt:       jwt,
		encryptor: encryptor,
		fs:        fs,
	}
}

func (s *ApplicationService) Create(c echo.Context) error {
	job := c.FormValue("job")
	if job == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Missing required fields",
			Error:   "Missing required fields",
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	application := new(models.ApplicationModel)
	application.ApplicantId = userId
	application.Job = job

	applicationId, err := s.store.Create(application)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "Application already exists",
				Error:   "Application already exists",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating application",
			Error:   err.Error(),
		})
	}

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

		switch file.Filename {
		case "policeClearance":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.POLICE_CLEARANCE_DIR,
			}
			url, err := s.fs.SaveFile(fileData)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving police clearance",
					Error:   err.Error(),
				})
			}

			clearance := new(models.PoliceClearanceModel)
			clearance.ApplicantId = userId
			clearance.ApplicationId = applicationId
			clearance.Url = url

			if _, err := s.store.NewPoliceClearance(clearance); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving police clerance",
					Error:   err.Error(),
				})
			}

		case "supportingDocument":
			fileData := fs.File{
				Data:     cipher,
				Category: fs.APPLICATION_PROOF_DIR,
			}
			url, err := s.fs.SaveFile(fileData)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving file",
					Error:   err.Error(),
				})
			}

			proof := new(models.ApplicationProofModel)
			proof.ApplicantId = userId
			proof.ApplicationId = applicationId
			proof.Url = url

			if _, err := s.store.NewProof(proof); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
					Message: "Error saving application proof",
					Error:   err.Error(),
				})
			}

		default:
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "File indicated by filename is unknown or unsupported",
				Error:   "Unsupported file",
			})

		}

	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"application": applicationId,
	})
}
