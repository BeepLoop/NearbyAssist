package response

import "nearbyassist/internal/models"

type VendorServiceList struct {
	Vendor struct {
		Id           string   `json:"id"`
		Name         string   `json:"name"`
		Email        string   `json:"email"`
		Phone        string   `json:"phone"`
		ImageUrl     string   `json:"imageUrl"`
		Rating       string   `json:"rating"`
		IsRestricted bool     `json:"isRestricted"`
		Expertise    []string `json:"expertise"`
		Socials      []string `json:"socials"`
	} `json:"vendor"`
	Services []struct {
		Id          string             `json:"id"`
		Title       string             `json:"title"`
		Description string             `json:"description"`
		Price       string             `json:"price"`
		Latitude    float64            `json:"latitude"`
		Longitude   float64            `json:"longitude"`
		Tags        []*models.TagModel `json:"tags"`
	} `json:"services"`
}
