package models

type ReviewModel struct {
	Model
	UpdateableModel
	ServiceId int    `json:"serviceId" db:"serviceId"`
	Rating    string `json:"rating" db:"rating"`
}

func (r *ReviewModel) Create() (int, error) {
	return 0, nil
}

func (r *ReviewModel) Update(id int) error {
	return nil
}

func (r *ReviewModel) Delete(id int) error {
	return nil
}
