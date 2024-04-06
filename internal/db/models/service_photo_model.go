package models

type ServicePhotoModel struct {
	Model
	UpdateableModel
	ServiceId int    `json:"serviceId" db:"serviceId"`
	VendorId  int    `json:"vendorId" db:"vendorId"`
	Url       string `json:"url" db:"url"`
}

func (s *ServicePhotoModel) Create() (int, error) {
	return 0, nil
}

func (s *ServicePhotoModel) Update(id int) error {
	return nil
}

func (s *ServicePhotoModel) Delete(id int) error {
	return nil
}
