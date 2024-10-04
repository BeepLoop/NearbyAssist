package response

type IdentityVerification struct {
	Id        string `db:"id" json:"id"`
	UserId    string `db:"userId" json:"userId"`
	Status    string `db:"status" json:"status"`
	CreatedAt string `db:"createdAt" json:"createdAt"`
}
