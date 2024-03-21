package types

type Acquaintance struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	ImageUrl string `db:"imageUrl"`
}
