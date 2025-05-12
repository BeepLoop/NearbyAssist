package dto

type UserAccountDetail struct {
	User           User
	ActiveBookings []Booking
	History        []Booking
}

type VendorAccountDetail struct {
	Vendor         Vendor
	Services       []Service
	ActiveBookings []Booking
	History        []Booking
}

type User struct {
	Id                         string
	Name                       string
	Email                      string
	ImageURL                   string
	Address                    string
	Phone                      string
	Socials                    []Social
	Identification             Identification
	CreatedAt                  string
	DateVerified               string
	IsRestricted               bool
	IsBanned                   bool
	HasSubmittedIdentification bool
}

type Vendor struct {
	Id             string
	Name           string
	Email          string
	ImageURL       string
	Address        string
	Phone          string
	Socials        []Social
	Identification Identification
	Rating         string
	Expertise      []Expertise
	JoinedAt       string
	DateVerified   string
	IsRestricted   bool
	IsBanned       bool
}

type Booking struct {
	Id            string
	Client        User
	Vendor        Vendor
	Service       Service
	Quantity      int
	Cost          string
	Status        string
	Extras        []Extra
	CancelReason  string
	CreatedAt     string
	ScheduleStart string
	ScheduleEnd   string
	UpdatedAt     string
}

type Extra struct {
	Id          string
	Title       string
	Description string
	Price       string
}

type Service struct {
	Id          string
	VendorId    string
	Title       string
	Description string
	Price       string
	PricingType string
	Tags        []string
	Extras      []Extra
	Images      []Image
	Address     Address
	CreatedAt   string
	UpdatedAt   string
}

type Image struct {
	Id  string
	URL string
}

type Expertise struct {
	Title              string
	DateApplied        string
	DateApproved       string
	SupportingDocument string
}

type Identification struct {
	Type          string
	IdNumber      string
	FrontImageURL string
	BackImageURL  string
}

type Address struct {
	Address   string
	Latitude  float64
	Longitude float64
}

type Social struct {
	Id    string
	Site  string
	Title string
	URL   string
}
