package types

type ServiceDetails struct {
	Id          int     `query:"id" db:"id"`
	VendorId    int     `query:"vendorId" db:"vendorId"`
	Title       string  `query:"title" db:"title"`
	Description string  `query:"description" db:"description"`
	VendorName  string  `query:"name" db:"name"`
	VendorImage string  `query:"imageUrl" db:"imageUrl"`
	Rate        string  `query:"rate" db:"rate"`
	Rating      float32 `query:"rating" db:"rating"`
	VendorRole  string  `query:"vendorRole" db:"vendorRole"`
	ReviewCount Count
}
