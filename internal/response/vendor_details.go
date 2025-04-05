package response

type VendorServices struct {
	Vendor  Vendor    `json:"vendor"`
	Servics []Service `json:"services"`
}
