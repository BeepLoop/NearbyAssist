package handlers

import (
	"nearbyassist/internal/encryption"
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

	req := models.NewIdentityVerificationModel(h.server.IdGen, h.server.DB)
	if req == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	} else {
		req.Name = name
		req.Address = address
		req.IdType = idType
		req.IdNumber = idNumber
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse submitted files")
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)

		switch file.Filename {
		case "frontId":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveFrontId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save front id")
			}

			frontId := models.NewFrontIdModelWithImageUrl(url, h.server.IdGen, h.server.DB)
			if frontId == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
			}

			if _, err := frontId.Create(); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save front id to db")
			} else {
				req.FrontId = frontId.Id
			}

		case "backId":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveBackId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save back id")
			}

			backId := models.NewBackIdModelWithImageUrl(url, h.server.IdGen, h.server.DB)
			if backId == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
			}

			if _, err := backId.Create(); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save back id to db")
			} else {
				req.BackId = backId.Id
			}

		case "face":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveFace)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save face")
			}

			face := models.NewFaceModelWithImageUrl(url, h.server.IdGen, h.server.DB)
			if face == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
			}

			if _, err := face.Create(); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save face to db")
			} else {
				req.Face = face.Id
			}

		}
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := req.EncryptName(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	if _, err := req.EncryptAddress(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	if _, err := req.EncryptIdNumber(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	if _, err := req.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while creating identity verify")
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message":        "Identity verification submitted.",
		"verificationId": req.Id,
	})
}

func (h *verificationHandler) HandleGetAllIdentityVerification(c echo.Context) error {
	ident_verification := models.NewIdentityVerificationModel(h.server.IdGen, h.server.DB)
	if ident_verification == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	requests, err := ident_verification.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while retrieving verification requests")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"requests": requests,
	})
}

func (h *verificationHandler) HandleGetIdentityVerification(c echo.Context) error {
	verificationId := c.Param("verificationId")
	if verificationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Verification ID must be a number")
	}

	ident_verification := models.NewIdentityVerificationModel(h.server.IdGen, h.server.DB)
	if ident_verification == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := ident_verification.FindById(verificationId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Verification request not found")
	}

	if _, err := ident_verification.DecryptName(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
	}

	if _, err := ident_verification.DecryptAddress(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
	}

	if _, err := ident_verification.DecryptIdNumber(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"request": ident_verification,
	})
}
