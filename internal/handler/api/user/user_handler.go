package user

import (
	"database/sql"
	"encoding/json"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	user_service "nearbyassist/internal/service/user"
	verification_service "nearbyassist/internal/service/verification"
	"nearbyassist/internal/utils"
	"net/http"
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
	bearerToken := utils.BearerTokenFromHeader(c)

	user, err := h.userService.GetUser(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error while retrieving user data",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"user": user,
	})
}

func (h *userHandler) FindUser(c echo.Context) error {
	userId := c.Param("userId")

	user, err := h.userService.GetUserById(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error while retrieving user data",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"user": user,
	})
}

func (h *userHandler) FindUserByEmail(c echo.Context) error {
	user, err := h.userService.FindByEmail(c.Param("email"))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "Email not found",
				Error:   "Email not found",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error while retrieving user data",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"user": user,
	})
}

func (h *userHandler) VerifyAccount(c echo.Context) error {
	user := c.FormValue("user")
	req := new(request.VerifyAccountPayload)
	if err := json.Unmarshal([]byte(user), req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.Error{
			Message: "Invalid payload",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error parsing form",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)
	if err := h.userVerificationService.VerifyAccount(bearerToken, req, files); err != nil {
		if strings.Contains(err.Error(), verification_service.ERR_ALREADY_VERIFIED) {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "Account already verified",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error processing verification request",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *userHandler) CheckVerificationStatus(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

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

	bearerToken := utils.BearerTokenFromHeader(c)
	socialId, err := h.userService.AddSocial(bearerToken, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error adding social",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"id": socialId,
	})
}

func (h *userHandler) DeleteSocial(c echo.Context) error {
	id := c.Param("id")
	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.userService.DeleteSocial(bearerToken, id); err != nil {
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

func (h *userHandler) AddExpertise(c echo.Context) error {
	expertiseId := c.FormValue("expertiseId")
	if expertiseId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Missing required fields",
			Error:   "Missing required fields",
		})
	}

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error reading files from form",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)
	if err := h.userService.AddVendorExpertise(bearerToken, expertiseId, files[0]); err != nil {
		if strings.Contains(err.Error(), user_service.ERR_FORBIDDEN) {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "Adding expertise not allowed",
				Error:   err.Error(),
			})
		}

		if strings.Contains(err.Error(), user_service.ERR_DUPLICATE_EXPERTISE) {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "You already have the applied expertise",
				Error:   err.Error(),
			})
		}

		if strings.Contains(err.Error(), user_service.ERR_DUPLICATE_APPLICATION) {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "You already have a pending application for the expertise",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error adding expertise",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *userHandler) SetDBL(c echo.Context) error {
	value := c.Param("value")

	bearerToken := utils.BearerTokenFromHeader(c)
	if err := h.userService.SetDBL(bearerToken, value); err != nil {
		if strings.Contains(err.Error(), user_service.ERR_INVALID_DBL) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, models.Error{
				Message: "Provided Daily Booking Limit value invalid",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error encountered while setting Daily Booking Limit",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *userHandler) ChangeAddress(c echo.Context) error {
	req := new(request.ChangeAddressPayload)
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

	bearerToken := utils.BearerTokenFromHeader(c)
	if err := h.userService.ChangeAddress(bearerToken, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error updating address",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
