package response

type TransactionSummary struct {
	Id                   string `db:"id"`
	CreatedAt            string `db:"createdAt"`
	Vendor               string `db:"vendor"`
	Client               string `db:"client"`
	ServiceTitle         string `db:"serviceTitle"`
	Price                string `db:"price"`
	StartDate            string `db:"startDate"`
	EndDate              string `db:"endDate"`
	VendorEmail          string `db:"vendorEmail"`
	ClientEmail          string `db:"clientEmail"`
	ConfirmationEndpoint string
}

type BasicEmailPayload struct {
	User            string
	SupportEndpoint string
}
