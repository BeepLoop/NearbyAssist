package admin

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/id_generator"
	store "nearbyassist/internal/store/admin"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminRefresh(t *testing.T) {
	adminStore := store.NewMockAdminStore()
	jwt := authenticator.NewMockAuthenticator()
	idGen := id_generator.NewMockGenerator()
	handler := NewHandler(adminStore, jwt, idGen)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/refresh", nil)

	router := echo.New()

	rec := httptest.NewRecorder()
	c := router.NewContext(req, rec)

	assert.NoError(t, handler.Refresh(c), "Should not return error")
	assert.Equal(t, http.StatusOK, rec.Code, "Should return status code 200")
}
