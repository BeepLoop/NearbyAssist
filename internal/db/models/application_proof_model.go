package models

type ApplicationProof struct {
	Model
	UpdateableModel
	ApplicationId int    `json:"applicationId" db:"applicationId"`
	ApplicantId   int    `json:"applicantId" db:"applicantId"`
	Url           string `json:"url" db:"url"`
}

func (a *ApplicationProof) Create() (int, error) {
	return 0, nil
}

func (a *ApplicationProof) Update(id int) error {
	return nil
}

func (a *ApplicationProof) Delete(id int) error {
	return nil
}
