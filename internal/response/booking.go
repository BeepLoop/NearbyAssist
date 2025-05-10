package response

type Booking struct {
	Id           string          `json:"id"`
	Vendor       User            `json:"vendor"`
	Client       User            `json:"client"`
	Service      ServiceBareInfo `json:"service"`
	Extras       []Extra         `json:"extras"`
	Cost         string          `json:"cost"`
	Status       string          `json:"status"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`
	ScheduledAt  string          `json:"scheduledAt"`
	CancelReason string          `json:"cancelReason"`
	CancelledBy  string          `json:"cancelledById"`
	QRSignature  string          `json:"qrSignature"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ServiceBareInfo struct {
	Id          string   `json:"id"`
	VendorId    string   `json:"vendorId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Rate        string   `json:"rate"`
	Tags        []Tag    `json:"tags"`
	Location    Location `json:"location"`
}

type Service struct {
	Id          string   `json:"id"`
	VendorId    string   `json:"vendorId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Rate        string   `json:"rate"`
	Tags        []Tag    `json:"tags"`
	Extras      []Extra  `json:"extras"`
	Images      []Image  `json:"images"`
	Location    Location `json:"location"`
	Disabled    bool     `json:"disabled"`
}

type Extra struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Image struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
