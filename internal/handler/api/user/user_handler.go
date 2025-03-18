package user

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	user_service "nearbyassist/internal/service/user"
	verification_service "nearbyassist/internal/service/verification"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userService             *user_service.Service
	userVerificationService *verification_service.Service
}

func NewHandler(userService *user_service.Service, userVerificationService *verification_service.Service) *userHandler {
	return &userHandler{
		userService:             userService,
		userVerificationService: userVerificationService,
	}
}

func (h *userHandler) GetUser(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	user, err := h.userService.GetUser(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.Error{
			Message: "Error while retrieving user data",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"user": user,
	})
}

func (h *userHandler) VerifyUserIdentity(c echo.Context) error {
	name := c.FormValue("name")
	phone := c.FormValue("phone")
	address := c.FormValue("address")
	latitude := c.FormValue("latitude")
	longitude := c.FormValue("longitude")
	idType := c.FormValue("idType")
	idNumber := c.FormValue("idNumber")
	if name == "" || phone == "" || address == "" || latitude == "" || longitude == "" || idType == "" || idNumber == "" {
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

	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid latitude",
			Error:   err.Error(),
		})
	}

	long, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid longitude",
			Error:   err.Error(),
		})
	}

	verificationId, err := h.userVerificationService.CreateVerificationRequest(
		name,
		phone,
		address,
		idType,
		idNumber,
		bearerToken,
		lat,
		long,
		files,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "Verification request already exists",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error processing verification request",
			Error:   err.Error(),
		})
	}

	user, err := h.userService.GetUser(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving user information",
			Error:   err.Error(),
		})
	}

	// TODO: Handle notifying the user of the verification status
	_ = user

	return c.JSON(http.StatusCreated, utils.Mapper{
		"verification": verificationId,
	})
}

func (h *userHandler) GetUserVerification(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	isVerified, err := h.userService.IsVerified(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Error checking if user is verified",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"verified": isVerified,
	})
}

func (h *userHandler) AddSocial(c echo.Context) error {
	req := new(request.AddSocialPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	if err := h.userService.AddSocial(bearerToken, req.Url); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error adding social",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h *userHandler) DeleteSocial(c echo.Context) error {
	req := new(request.DeleteSocialPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	if err := h.userService.DeleteSocial(bearerToken, req.Url); err != nil {
		if strings.Contains(err.Error(), "social not found") {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "social does not exists",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error deleting social",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
