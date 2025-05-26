package dto

type PendingService struct {
	Vendor  Vendor
	Service Service
}

type PendingServiceListItem struct {
	ID          string
	VendorName  string
	VendorEmail string
	CreatedAt   string
}
