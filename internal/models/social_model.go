package models

type SocialModel struct {
	Model
	UserId string `db:"userId"`
	Site   string `db:"site"`
	Title  string `db:"title"`
	Url    string `db:"url"`
}
