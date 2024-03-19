package types

type Session struct {
	Id     int    `db:"id"`
	UserId int    `db:"userId"`
	Token  string `db:"token"`
}
