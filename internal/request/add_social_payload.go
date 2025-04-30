package request

type AddSocialPayload struct {
	Site  string `json:"site" validate:"required"`
	Title string `json:"title" validate:"required"`
	Url   string `json:"url" validate:"required"`
}
