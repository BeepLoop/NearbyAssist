package request

type AddExtraPayload struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	ServiceId   string  `json:"serviceId" validate:"required"`
}
