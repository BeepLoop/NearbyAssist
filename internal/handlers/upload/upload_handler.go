package upload

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/upload/v1"

	"github.com/labstack/echo/v4"
)

func UploadHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	r.PUT("/service", upload.ServicePhoto).Name = "upload service image"
	r.PUT("/application/proof", upload.VendorApplicationProof).Name = "upload vendor application proof"
}
