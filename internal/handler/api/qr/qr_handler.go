package qr

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	qr_service "nearbyassist/internal/service/qr"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type qrHandler struct {
	service *qr_service.Service
}

func NewHandler(service *qr_service.Service) *qrHandler {
	return &qrHandler{service: service}
}

func (h *qrHandler) SignTransaction(c echo.Context) error {
	req := new(request.QRSignatureInput)
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

	signature, err := h.service.SignData(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error occurred while generating signature",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"signature": signature,
	})
}

func (h *qrHandler) VerifySignature(c echo.Context) error {
	req := new(request.QRSignatureVerifyInput)
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

	ok := h.service.VerifySignature(req)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Invalid signature",
			Error:   "Invalid signature",
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
