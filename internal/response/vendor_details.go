package response

import "nearbyassist/internal/models"

type DetailVendorResponse struct {
	Vendor   *models.VendorModel    `json:"vendor"`
	Services []*models.ServiceModel `json:"service"`
}
