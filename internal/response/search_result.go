package response

type SearchResult struct {
	Id        string  `json:"id"`
	Score     float32 `json:"score"`
	Rank      int     `json:"rank"`
	Vendor    string  `json:"vendor"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
