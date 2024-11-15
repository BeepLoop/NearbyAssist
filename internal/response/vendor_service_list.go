package response

type VendorServiceList struct {
	Vendor struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		ImageUrl     string `json:"imageUrl"`
		Rating       string `json:"rating"`
		IsRestricted int    `json:"isRestricted"`
	} `json:"vendor"`
	Services []struct {
		Id          string   `json:"id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       string   `json:"price"`
		Latitude    float64  `json:"latitude"`
		Longitude   float64  `json:"longitude"`
		Tags        []string `json:"tags"`
	} `json:"services"`
}
