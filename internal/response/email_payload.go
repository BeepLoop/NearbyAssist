package response

type TransactionSummary struct {
	Id                   string
	CreatedAt            string
	ServiceProvider      string
	Client               string
	ServiceTitle         string
	ServiceCategory      string
	Price                string
	StartDate            string
	Location             string
	ProviderEmail        string
	ClientEmail          string
	ConfirmationEndpoint string
}

type BasicEmailPayload struct {
	User            string
	SupportEndpoint string
}
