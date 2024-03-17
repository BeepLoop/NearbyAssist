package types

type Service struct {
	Id          int     `db:"id"`
	Vendor      int     `query:"vendorId" db:"vendorId"`
	Title       string  `query:"title" db:"title"`
	Description string  `query:"description" db:"description"`
	Rate        float64 `query:"rate" db:"rate"`
	Location    string  `query:"location" db:"location"`
	Category    string  `query:"category" db:"category"`
}

type TransformedServiceData struct {
	Id          int
	Vendor      int
	Title       string
	Description string
	Rate        float64
	Latitude    float64
	Longitude   float64
}
