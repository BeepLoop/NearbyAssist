package response

type Booking struct {
	Id                 string         `json:"id"`
	Vendor             User           `json:"vendor"`
	Client             User           `json:"client"`
	ServiceId          string         `json:"serviceId"`
	ServiceTitle       string         `json:"serviceTitle"`
	ServiceDescription string         `json:"serviceDescription"`
	Price              string         `json:"price"`
	PricingType        string         `json:"pricingType"`
	Extras             []BookingExtra `json:"extras"`
	Quantity           int            `json:"quantity"`
	Cost               string         `json:"cost"`
	Status             string         `json:"status"`
	CreatedAt          string         `json:"createdAt"`
	UpdatedAt          string         `json:"updatedAt"`
	RequestedStart     string         `json:"requestedStart"`
	RequestedEnd       string         `json:"requestedEnd"`
	ScheduleStart      string         `json:"scheduleStart"`
	ScheduleEnd        string         `json:"scheduleEnd"`
	CancelReason       string         `json:"cancelReason"`
	CancelledBy        string         `json:"cancelledById"`
	QRSignature        string         `json:"qrSignature"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

type ServiceBareInfo struct {
	Id          string   `json:"id"`
	VendorId    string   `json:"vendorId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       string   `json:"price"`
	PricingType string   `json:"pricingType"`
	Tags        []Tag    `json:"tags"`
	Location    Location `json:"location"`
}

type Service struct {
	Id          string   `json:"id"`
	VendorId    string   `json:"vendorId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       string   `json:"price"`
	PricingType string   `json:"pricingType"`
	Tags        []Tag    `json:"tags"`
	Extras      []Extra  `json:"extras"`
	Images      []Image  `json:"images"`
	Location    Location `json:"location"`
	Disabled    bool     `json:"disabled"`
}

type BookingExtra struct {
	BookingId   string `json:"bookingId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type Extra struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type Image struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
