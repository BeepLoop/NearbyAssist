package request

type AddSocialPayload struct {
	Url string `json:"url" validate:"required"`
}
