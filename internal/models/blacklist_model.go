package models

type BlacklistModel struct {
	Id    int    `json:"id" db:"id"`
	Token string `json:"token" db:"token"`
}

func NewBlacklistModel(token string) *BlacklistModel {
	return &BlacklistModel{
		Token: token,
	}
}
