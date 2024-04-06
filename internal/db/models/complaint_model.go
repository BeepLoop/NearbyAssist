package models

type ComplaintModel struct {
	Model
	UpdateableModel
	Code    int    `json:"code" db:"code"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

func (c *ComplaintModel) Create() (int, error) {
	return 0, nil
}

func (c *ComplaintModel) Update(id int) error {
	return nil
}

func (c *ComplaintModel) Delete(id int) error {
	return nil
}
