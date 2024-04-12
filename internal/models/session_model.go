package models

type SessionModel struct {
	Model
	UpdateableModel
	UserId int    `json:"userId" db:"userId"`
	Status string `json:"status" db:"status"`
	Token  string `json:"token" db:"token"`
}

func NewSessionModel(userId int, token string) *SessionModel {
	return &SessionModel{
		UserId: userId,
		Token:  token,
	}
}
