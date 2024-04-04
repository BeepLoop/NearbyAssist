package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetUploadParams(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/v1/uploads/upload?vendorId=1&serviceId=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	params, err := GetUploadParams(c, "vendorId", "serviceId")
	if err != nil {
		t.Fatalf("Expected nil but got %v", err)
	}

	expected := map[string]int{
		"vendorId":  1,
		"serviceId": 2,
	}

	assert.EqualValues(t, params, expected)
}
