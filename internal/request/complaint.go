package request

type SystemComplaintPayload struct {
	Title  string `json:"title" validate:"required"`
	Detail string `json:"detail" validate:"required"`
}
