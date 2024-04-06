package models

type VendorModel struct {
	Model
	UpdateableModel
	VendorId   int    `json:"vendorId" db:"vendorId"`
	Rating     string `json:"rating" db:"rating"`
	Job        string `json:"job" db:"job"`
	Restricted int    `json:"restricted" db:"restricted"`
}

func (v *VendorModel) Create() (int, error) {
	return 0, nil
}

func (v *VendorModel) Update(id int) error {
	return nil
}

func (v *VendorModel) Delete(id int) error {
	return nil
}
