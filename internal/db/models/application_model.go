package models

type ApplicationModel struct {
	Model
	UpdateableModel
    GeoSpatialModel
	ApplicantId int    `json:"applicantId" db:"applicantId"`
	Job         string `json:"job" db:"job"`
	Status      string `json:"status" db:"status"`
}

func (a *ApplicationModel) Create() (int, error) {
	return 0, nil
}

func (a *ApplicationModel) Update(id int) error {
	return nil
}

func (a *ApplicationModel) Delete(id int) error {
	return nil
}
