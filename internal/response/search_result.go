package response

type ServiceSearchResult struct {
	Id                    string  `json:"id"`
	VendorName            string  `json:"vendorName"`
	SuggestionScore       float32 `json:"suggestionScore"`
	Rate                  float32 `json:"rate"`
	Rating                float32 `json:"rating"`
	Latitude              float64 `json:"latitude"`
	Longitude             float64 `json:"longitude"`
	CompletedTransactions float32 `json:"completedTransactions"`
	Distance              float32 `json:"distance"`
}
