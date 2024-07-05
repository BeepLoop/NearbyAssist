package response

type AllVerification struct {
	Id        int    `json:"id" db:"id"`
	User      int    `json:"user" db:"user"`
	CreatedAt string `json:"createdAt" db:"createdAt"`
}
