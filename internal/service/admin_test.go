package service

import (
	"nearbyassist/internal/store/admin"
	"nearbyassist/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminLogin(t *testing.T) {
	adminStore := admin.NewMockAdminStore()
	handler := NewAdminService(adminStore)

	t.Run("Should fail if payload is invalid", func(t *testing.T) {
		tests := []struct {
			payload      string
			expectedCode int
		}{
			{
				`{"username": "foo", "password": ""}`,
				http.StatusBadRequest,
			},
			{
				`{"username": "foo"}`,
				http.StatusBadRequest,
			},
			{
				`{"username": "foo", "": "bar"}`,
				http.StatusBadRequest,
			},
			{
				`{"username": "foo", "": "bar"`,
				http.StatusBadRequest,
			},
			{
				`{"Username": "foo", "Password": ""}`,
				http.StatusBadRequest,
			},
			{
				`"Username": "foo", "Password": ""`,
				http.StatusBadRequest,
			},
		}

		for _, test := range tests {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/login", strings.NewReader(test.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			router := echo.New()
			router.Validator = &utils.Validator{Validator: validator.New()}
			c := router.NewContext(req, rec)

			assert.Error(t, handler.Login(c), "Should return error")
		}
	})
}

func TestAdminRefresh(t *testing.T) {
	adminStore := admin.NewMockAdminStore()
	handler := NewAdminService(adminStore)

	t.Run("Should fail if payload is invalid", func(t *testing.T) {
		tests := []struct {
			payload      string
			expectedCode int
		}{
			{
				`{}`,
				http.StatusBadRequest,
			},
			{
				`{"Token": ""}`,
				http.StatusBadRequest,
			},
		}

		for _, test := range tests {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/refresh", strings.NewReader(test.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			router := echo.New()
			router.Validator = &utils.Validator{Validator: validator.New()}
			c := router.NewContext(req, rec)

			assert.Error(t, handler.Refresh(c), "Should return error")
		}
	})
}
