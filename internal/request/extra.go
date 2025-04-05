package request

type NewExtra struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}

type Extra struct {
	Id          string  `json:"id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}
