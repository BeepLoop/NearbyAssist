package models

type MessageModel struct {
	Model
	UpdateableModel
	Sender   int    `json:"sender" db:"sender"`
	Receiver int    `json:"receiver" db:"receiver"`
	Content  string `json:"content" db:"content"`
}

func (m *MessageModel) Create() (int, error) {
	return 0, nil
}

func (m *MessageModel) Update(id int) error {
	return nil
}

func (m *MessageModel) Delete(id int) error {
	return nil
}
