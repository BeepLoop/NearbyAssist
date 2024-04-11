package models

// import "nearbyassist/internal/db"
type SessionModel struct {
	Model
	UpdateableModel
	UserId int    `json:"userId" db:"userId"`
	Status string `json:"status" db:"status"`
	Token  string `json:"token" db:"token"`
}

func NewSessionModel() *SessionModel {
	return &SessionModel{}
}

//
// func (s *SessionModel) Create() (int, error) {
// 	return 0, nil
// }
//
// func (s *SessionModel) Update(id int) error {
// 	return nil
// }
