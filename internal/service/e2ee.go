package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/store/e2ee"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type E2EEService struct {
	store   e2ee.E2EEStore
	jwt     auth.Authenticator
	encrypt auth.Encryption
}

func NewE2EEService(store e2ee.E2EEStore, jwt auth.Authenticator, encrypt auth.Encryption) *E2EEService {
	return &E2EEService{
		store:   store,
		jwt:     jwt,
		encrypt: encrypt,
	}
}

func (s *E2EEService) SaveKeys(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	req := new(request.PEM)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
	}

	publicKey := new(models.PublicKeyModel)
	publicKey.Owner = userId
	publicKey.Pem = req.Public
	if _, err := s.store.NewPublicPem(publicKey); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error saving public key",
			Error:   err.Error(),
		})
	}

	privateKey := new(models.PrivateKeyModel)
	privateKey.Owner = userId
	if encrypted, err := s.encrypt.EncryptString(req.Private); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ENCRYPTION_ERR,
			Error:   err.Error(),
		})
	} else {
		privateKey.Pem = encrypted
	}

	if _, err := s.store.NewPrivatePem(privateKey); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error saving private key",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (s *E2EEService) GetKeys(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	private, err := s.store.GetPrivatePem(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Private key not found",
			Error:   err.Error(),
		})
	}

	decrypted, err := s.encrypt.DecryptString(private.Pem)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	}

	public, err := s.store.GetPublicPem(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Public key not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"publicKey":  public.Pem,
		"privateKey": decrypted,
	})
}

func (s *E2EEService) GetPublicKey(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "userId must be an integer",
			Error:   "userId must be an integer",
		})
	}

	publicKey, err := s.store.GetPublicPem(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Public key not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"key": publicKey.Pem,
	})
}
