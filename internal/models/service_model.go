package models

type ServiceSearchResult struct {
	ServiceModel
	Vendor string `json:"vendor" db:"vendor"`
}

type PricingType string

const (
	FIXED_PRICING  PricingType = "fixed"
	HOURLY_PRICING PricingType = "per_hour"
	DAILY_PRICING  PricingType = "per_day"
)

type ServiceModel struct {
	Model
	UpdateableModel
	VendorId    string      `db:"vendorId"`
	Title       string      `db:"title"`
	Description string      `db:"description"`
	Price       string      `db:"price"`
	PricingType PricingType `db:"pricingType"`
	Signature   string      `db:"signature"`
	Disabled    bool        `db:"disabled"`
	Address     AddressModel

	// Additional fields for joins
	Vendor       VendorModel
	Tags         []*TagModel          `json:"tags" db:"tags" validate:"required"`
	TagsAsString []string             `json:"-"`
	Extras       []*ExtraModel        `json:"extras" db:"extras"`
	Images       []*ServicePhotoModel `json:"images"`
}
