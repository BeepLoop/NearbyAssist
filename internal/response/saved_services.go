package response

import "nearbyassist/internal/models"

type SavedServicesResponse struct {
	Services []SavedServiceData `json:"services"`
}

type SavedServiceData struct {
	Service *models.ServiceModel `json:"serviceInfo"`
	Vendor  struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		ImageUrl     string `json:"imageUrl"`
		Rating       string `json:"rating"`
		IsRestricted bool   `json:"isRestricted"`
	} `json:"vendorInfo"`
	Photos         []*models.ServicePhotoModel `json:"serviceImages"`
	CountPerRating CountPerRating              `json:"countPerRating"`
}
