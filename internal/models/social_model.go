package models

type SocialModel struct {
	Model
	UserId string `db:"userId"`
	Url    string `db:"url"`
}
