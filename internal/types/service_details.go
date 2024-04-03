package types

type ServiceDetails struct {
	Id          int    `query:"id" db:"id"`
	VendorId    int    `query:"vendorId" db:"vendorId"`
	Title       string `query:"title" db:"title"`
	Description string `query:"description" db:"description"`
	VendorName  string `query:"name" db:"name"`
	VendorImage string `query:"imageUrl" db:"imageUrl"`
	Rate        string `query:"rate" db:"rate"`
	Rating      string `query:"rating" db:"rating"`
	Job         string `query:"job" db:"job"`
	ReviewCount Count
	Photos      []Photo
}
