package types

type Review struct {
	Id        int `db:"id"`
	ServiceId int `db:"serviceId"`
	Rating    int `db:"rating"`
}

type ReviewCount struct {
	Rating string `db:"rating"`
	Count  int    `db:"count"`
}
