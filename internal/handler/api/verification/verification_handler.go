package verification

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/email"
	user_service "nearbyassist/internal/service/user"
	verification_service "nearbyassist/internal/service/verification"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type verificationHandler struct {
	verificationService *verification_service.Service
	userService         *user_service.Service
	mailman             email.MailService
}

func NewHandler(verificationService *verification_service.Service, userService *user_service.Service, mailman email.MailService) *verificationHandler {
	return &verificationHandler{verificationService: verificationService, userService: userService, mailman: mailman}
}

func (h *verificationHandler) CreateIdentityVerification(c echo.Context) error {
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

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	verificationId, err := h.verificationService.CreateVerificationRequest(
		name,
		address,
		idType,
		idNumber,
		bearerToken,
		files,
	)

	user, err := h.userService.GetUser(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving user information",
			Error:   err.Error(),
		})
	}

	go email.IdentityVerificationMail(h.mailman).To([]string{user.Email}).Send()

	return c.JSON(http.StatusCreated, utils.Mapper{
		"verification": verificationId,
	})
}
