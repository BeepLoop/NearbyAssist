package handlers

import (
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type verificationHandler struct {
	server *server.Server
}

func NewVerificationHandler(s *server.Server) *verificationHandler {
	return &verificationHandler{
		server: s,
	}
}

func (h *verificationHandler) HandleVerifyIdentity(c echo.Context) error {
	name := c.FormValue("name")
	address := c.FormValue("address")
	idType := c.FormValue("idType")
	idNumber := c.FormValue("idNumber")

	req := &models.IdentityVerificationModel{
		Name:     name,
		Address:  address,
		IdType:   idType,
		IdNumber: idNumber,
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse submitted files")
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)

		switch file.Filename {
		case "frontId":
			filename, err := handler.SavePhoto(file, h.server.Storage.SaveFrontId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save front id")
			}

			if id, err := h.server.DB.NewFrontId(&models.FrontIdModel{Url: filename}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save front id to db")
			} else {
				req.FrontId = id
			}
		case "backId":
			filename, err := handler.SavePhoto(file, h.server.Storage.SaveBackId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save back id")
			}

			if id, err := h.server.DB.NewBackId(&models.BackIdModel{Url: filename}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save back id to db")
			} else {
				req.BackId = id
			}
		case "face":
			filename, err := handler.SavePhoto(file, h.server.Storage.SaveFace)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save face")
			}

			if id, err := h.server.DB.NewFace(&models.FaceModel{Url: filename}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save face to db")
			} else {
				req.Face = id
			}
		}
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	verificationId, err := h.server.DB.NewIdentityVerification(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message":        "Identity verification submitted.",
		"verificationId": verificationId,
	})
}
