package response

type Review struct {
	Id               string `json:"id"`
	BookingId        string `json:"bookingId"`
	Rating           int    `json:"rating"`
	Text             string `json:"text"`
	CreatedAt        string `json:"createdAt"`
	RevieweeName     string `json:"revieweeName"`
	RevieweeImageUrl string `json:"revieweeImageUrl"`
}
