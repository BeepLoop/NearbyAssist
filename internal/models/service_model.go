package models

type ServiceSearchResult struct {
	ServiceModel
	Vendor string `json:"vendor" db:"vendor"`
}

type ServiceModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	VendorId    string `json:"vendorId" db:"vendorId" validate:"required"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description" validate:"required"`
	Rate        string `json:"rate" db:"rate" validate:"required"`
	Signature   string `db:"signature" json:"-"`

	// Additional fields for joins
	Tags         []*TagModel          `json:"tags" db:"tags" validate:"required"`
	TagsAsString []string             `json:"-"`
	Extras       []*ExtraModel        `json:"extras" db:"extras"`
	Images       []*ServicePhotoModel `json:"images"`
}
