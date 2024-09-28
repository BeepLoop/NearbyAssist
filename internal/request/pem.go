package request

type PEM struct {
	Private string `json:"privatePem" validate:"required"`
	Public  string `json:"publicPem" validate:"required"`
}
