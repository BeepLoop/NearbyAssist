package admin

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/id_generator"
	store "nearbyassist/internal/store/admin"
	"nearbyassist/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminRefresh(t *testing.T) {
	adminStore := store.NewMockAdminStore()
	jwt := authenticator.NewMockAuthenticator()
	idGen := id_generator.NewMockGenerator()
	handler := NewHandler(adminStore, jwt, idGen)

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
