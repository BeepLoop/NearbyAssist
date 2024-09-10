package service

import (
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/store/user"
	"nearbyassist/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	userStore := user.NewMockUserStore()
	encryptor := auth.NewMockEncryptor()
	handler := NewUserService(userStore, encryptor)

	t.Run("Should fail if payload is invalid", func(t *testing.T) {
		tests := []struct {
			payload      string
			expectedCode int
		}{
			{
				`{"name": "", "email": "", "image": ""}`,
				http.StatusBadRequest,
			},
			{
				`{"name": "", "email": "", "imageUrl": ""}`,
				http.StatusBadRequest,
			},
		}

		for _, test := range tests {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", strings.NewReader(test.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			router := echo.New()
			router.Validator = &utils.Validator{Validator: validator.New()}
			c := router.NewContext(req, rec)

			assert.Error(t, handler.Login(c), "Should return error")
		}
	})
}
