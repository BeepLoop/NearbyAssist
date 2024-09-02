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
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error initializing model",
            Error:   models.MODEL_INIT_ERROR,
        })
	} else {
		req.Name = name
		req.Address = address
		req.IdType = idType
		req.IdNumber = idNumber
	}

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error parsing form",
            Error:   err.Error(),
        })
	}

	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)

		switch file.Filename {
		case "frontId":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveFrontId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error saving front ID",
                    Error:   err.Error(),
                })
			}

			frontId := models.NewFrontIdModelWithImageUrl(url, h.server.IdGen, h.server.DB)
			if frontId == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error initializing model",
                    Error:   models.MODEL_INIT_ERROR,
                })
			}

			if _, err := frontId.Create(); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error creating front ID",
                    Error:   err.Error(),
                })
			} else {
				req.FrontId = frontId.Id
			}

		case "backId":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveBackId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error saving back ID",
                    Error:   err.Error(),
                })
			}

			backId := models.NewBackIdModelWithImageUrl(url, h.server.IdGen, h.server.DB)
			if backId == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error initializing model",
                    Error:   models.MODEL_INIT_ERROR,
                })
			}

			if _, err := backId.Create(); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error creating back ID",
                    Error:   err.Error(),
                })
			} else {
				req.BackId = backId.Id
			}

		case "face":
			url, err := handler.SavePhoto(file, h.server.Storage.SaveFace)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error saving face",
                    Error:   err.Error(),
                })
			}

			face := models.NewFaceModelWithImageUrl(url, h.server.IdGen, h.server.DB)
			if face == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error initializing model",
                    Error:   models.MODEL_INIT_ERROR,
                })
			}

			if _, err := face.Create(); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
                    Message: "Error creating face",
                    Error:   err.Error(),
                })
			} else {
				req.Face = face.Id
			}

		}
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
            Message: "Error validating request body",
            Error:   err.Error(),
        })
	}

	if _, err := req.EncryptName(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error encrypting name",
            Error:   encryption.ENCRYPTION_ERR,
        })
	}

	if _, err := req.EncryptAddress(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error encrypting address",
            Error:   encryption.ENCRYPTION_ERR,
        })
	}

	if _, err := req.EncryptIdNumber(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error encrypting ID number",
            Error:   encryption.ENCRYPTION_ERR,
        })
	}

	if _, err := req.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error creating identity verification",
            Error:   err.Error(),
        })
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message":        "Identity verification submitted.",
		"verificationId": req.Id,
	})
}

func (h *verificationHandler) HandleGetAllIdentityVerification(c echo.Context) error {
	ident_verification := models.NewIdentityVerificationModel(h.server.IdGen, h.server.DB)
	if ident_verification == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error initializing model",
            Error:   models.MODEL_INIT_ERROR,
        })
	}

	params := utils.ParseQuery(c.QueryString())
	requests, err := ident_verification.FindAll(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error getting identity verifications",
            Error:   err.Error(),
        })
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"requests": requests,
	})
}

func (h *verificationHandler) HandleGetIdentityVerification(c echo.Context) error {
	verificationId := c.Param("verificationId")
	if verificationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
            Message: "Verification ID is required",
            Error:   "Verification ID is required",
        })
	}

	ident_verification := models.NewIdentityVerificationModel(h.server.IdGen, h.server.DB)
	if ident_verification == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error initializing model",
            Error:   models.MODEL_INIT_ERROR,
        })
	}

	if _, err := ident_verification.FindById(verificationId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
            Message: "Identity verification not found",
            Error:   err.Error(),
        })
	}

	if _, err := ident_verification.DecryptName(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error decrypting name",
            Error:   encryption.DECRYPTION_ERR,
        })
	}

	if _, err := ident_verification.DecryptAddress(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error decrypting address",
            Error:   encryption.DECRYPTION_ERR,
        })
	}

	if _, err := ident_verification.DecryptIdNumber(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
            Message: "Error decrypting ID number",
            Error:   encryption.DECRYPTION_ERR,
        })
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"request": ident_verification,
	})
}
