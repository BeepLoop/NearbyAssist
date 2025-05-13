package response

type ServiceSearchResult struct {
	Id                string  `json:"id"`
	VendorName        string  `json:"vendorName"`
	Suggestibility    float32 `json:"suggestibility"`
	Price             string  `json:"price"`
	Rating            string  `json:"rating"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	CompletedBookings int     `json:"completedBookings"`
	Distance          float32 `json:"distance"`
	Service           Service `json:"service"`
}
