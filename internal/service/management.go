package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/store/management"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ManagementService struct {
	store     management.ManagementStore
	jwt       auth.Authenticator
	encryptor auth.Encryption
	hash      auth.Hash
}

func NewManagementService(store management.ManagementStore, jwt auth.Authenticator, encryptor auth.Encryption, hash auth.Hash) *ManagementService {
	return &ManagementService{
		store:     store,
		jwt:       jwt,
		encryptor: encryptor,
		hash:      hash,
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

	usernameHash, err := s.hash.Generate([]byte(req.Username))
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

func (s *ManagementService) GetUser(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID must be a number",
			Error:   "User ID must be a number",
		})
	}

	user, err := s.store.GetUserById(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to get user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"user": user,
	})
}

func (s *ManagementService) RestrictVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	if err := s.store.RestrictVendor(vendorId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error restricting vendor",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (s *ManagementService) UnrestrictVendor(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Vendor ID must be a number",
			Error:   "Vendor ID must be a number",
		})
	}

	if err := s.store.UnrestrictVendor(vendorId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error restricting vendor",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (s *ManagementService) GetAllApplications(c echo.Context) error {
	params := utils.ParseQuery(c.QueryString())
	results, err := s.store.GetApplications(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting applications",
			Error:   err.Error(),
		})
	}

	applications := make([]response.ApplicationPayload, 0)
	for _, result := range results {
		applications = append(applications, response.ApplicationPayload{
			Id:          result.Id,
			ApplicantId: result.ApplicantId,
			Status:      string(result.Status),
			CreatedAt:   result.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"applications": applications,
	})
}

func (s *ManagementService) GetApplication(c echo.Context) error {
	applicationId := c.Param("applicationId")
	if applicationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Application ID is required",
			Error:   "Application ID is required",
		})
	}

	application, err := s.store.GetApplicationById(applicationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Application not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"application": application,
	})
}

func (s *ManagementService) ApproveApplication(c echo.Context) error {
	applicationId := c.Param("applicationId")
	if applicationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Application ID is required",
			Error:   "Application ID is required",
		})
	}

	if err := s.store.ApproveApplication(applicationId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error approving application",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (s *ManagementService) RejectApplication(c echo.Context) error {
	applicationId := c.Param("applicationId")
	if applicationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Application ID is required",
			Error:   "Application ID is required",
		})
	}

	if err := s.store.RejectApplication(applicationId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error approving application",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (s *ManagementService) GetTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Transaction ID is required",
			Error:   "Transaction ID is required",
		})
	}

	transaction, err := s.store.GetTransaction(transactionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting transaction",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transaction": transaction,
	})
}

func (s *ManagementService) GetSystemComplaints(c echo.Context) error {
	params := utils.ParseQuery(c.QueryString())
	complaints, err := s.store.GetSystemComplaints(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting system complaints",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaints": complaints,
	})
}

func (s *ManagementService) GetSystemComplaint(c echo.Context) error {
	complaintId := c.Param("complaintId")
	if complaintId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Complaint ID is required",
			Error:   "Complaint ID is required",
		})
	}

	complaint, err := s.store.GetSystemComplaintById(complaintId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting system complaint",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaint": complaint,
	})
}

func (s *ManagementService) GetVerificationRequests(c echo.Context) error {
	params := utils.ParseQuery(c.QueryString())
	requests, err := s.store.GetVerificationRequests(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting verification requests",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"requests": requests,
	})
}

func (s *ManagementService) HandleGetIdentityVerification(c echo.Context) error {
	verificationId := c.Param("verificationId")
	if verificationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Verification ID is required",
			Error:   "Verification ID is required",
		})
	}

	request, err := s.store.GetVerificationRequestById(verificationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting verification request",
			Error:   err.Error(),
		})
	}

	if decryptedName, err := s.encryptor.DecryptString(request.Name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		request.Name = decryptedName
	}

	if decryptedAddress, err := s.encryptor.DecryptString(request.Address); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		request.Address = decryptedAddress
	}

	if decryptedIdNumber, err := s.encryptor.DecryptString(request.IdNumber); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		request.IdNumber = decryptedIdNumber
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"request": request,
	})
}

func (s *ManagementService) ApproveIdentityVerification(c echo.Context) error {
	verificationId := c.Param("verificationId")
	if verificationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Verification ID is required",
			Error:   "Verification ID is required",
		})
	}

	if err := s.store.ApproveIdentityVerification(verificationId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error approving verification request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
