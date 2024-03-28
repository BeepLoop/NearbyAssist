package types

type TransactionHistory struct {
	Id        int    `db:"id"`
	Vendor    string `db:"vendor"`
	Client    string `db:"client"`
	Service   string `db:"service"`
	Status    string `db:"status"`
}
