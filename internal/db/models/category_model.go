package models

type CategoryModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`
}

func (c *CategoryModel) Create() (int, error) {
	return 0, nil
}

func (c *CategoryModel) Update(id int) error {
	return nil
}

func (c *CategoryModel) Delete(id int) error {
	return nil
}
