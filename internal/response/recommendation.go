package response

import "nearbyassist/internal/models"

type Recommendation struct {
	Searches []string                 `json:"searches"`
	Services []*ServiceRecommendation `json:"services"`
}

type ServiceRecommendation struct {
	Id          string             `json:"id"`
	VendorId    string             `json:"vendorId"`
	Vendor      string             `json:"vendor"`
	Thumbnail   string             `json:"thumbnail"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Rating      string             `json:"rating"`
	Rate        string             `json:"rate"`
	Tags        []*models.TagModel `json:"tags"`
}
