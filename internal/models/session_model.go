package models

type SessionModel struct {
	Model
	UpdateableModel
	Status string `json:"status" db:"status"`
	Token  string `json:"token" db:"token"`
}

func NewSessionModel(token string) *SessionModel {
	return &SessionModel{
		Token: token,
	}
}
