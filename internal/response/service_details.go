package response

type DetailedServiceResponse struct {
	Vendor  Vendor   `json:"vendor"`
	Service Service  `json:"service"`
	Ratings []int    `json:"ratings"` // Length of 5, position indicates level
	Reviews []Review `json:"reviews"`
}
