package types

type Service struct {
	Vendor      int     `query:"vendor" db:"vendor"`
	Title       string  `query:"title" db:"title"`
	Description string  `query:"description" db:"description"`
	Rate        float64 `query:"rate" db:"rate"`
	Location    string  `query:"location" db:"location"`
	Category    string  `query:"category" db:"category"`
}

type TransformedServiceData struct {
	Vendor      int
	Title       string
	Description string
	Rate        float64
	Latitude    float64
	Longitude   float64
}
