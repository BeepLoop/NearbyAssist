package response

type SearchResult struct {
	Id             int     `json:"id"`
	Suggestability float32 `json:"suggestability"`
	Rank           int     `json:"rank"`
	Vendor         string  `json:"vendor"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
}
