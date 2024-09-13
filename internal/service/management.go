package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/store/management"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ManagementService struct {
	store     management.ManagementStore
	encryptor auth.Encryption
}

func NewManagementService(store management.ManagementStore, encryptor auth.Encryption) *ManagementService {
	return &ManagementService{
		store:     store,
		encryptor: encryptor,
	}
}

func (s *ManagementService) CreateStaff(c echo.Context) error {
	req := new(request.NewAdminPayload)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
	}

	newAdmin := new(models.AdminModel)

	usernameHash, err := auth.Sha256([]byte(req.Username))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to hash username",
			Error:   err.Error(),
		})
	}

	encryptedUsername, err := s.encryptor.EncryptString(req.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to encrypt username",
			Error:   err.Error(),
		})
	}

	hashedPassword, err := auth.BcryptPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
	}

	newAdmin.Username = encryptedUsername
	newAdmin.UsernameHash = usernameHash
	newAdmin.Role = models.AdminRole(req.Role)
	newAdmin.Password = hashedPassword

	adminId, err := s.store.CreateStaff(newAdmin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to create staff",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"adminId": adminId,
	})
}
