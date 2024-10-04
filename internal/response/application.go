package response

type ApplicationPayload struct {
	Id          string `json:"id"`
	ApplicantId string `json:"applicantId"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}
