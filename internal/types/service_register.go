package types

type ServiceRegister struct {
	VendorId    int     `db:"vendorId" json:"vendorId" validate:"required"`
	Title       string  `db:"title" json:"title" validate:"required"`
	Description string  `db:"description" json:"description" validate:"required"`
	Rate        float64 `db:"rate" json:"rate" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
	Point       string  `db:"point"`
	CategoryId  int     `db:"categoryId" json:"categoryId" validate:"required"`
}
