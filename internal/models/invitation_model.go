package models

type InvitationModel struct {
	Id           string `db:"id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Code         string `db:"code"`
	Password     string `db:"password"`
	UsernameHash string `db:"usernameHash"`
	EmailHash    string `db:"emailHash"`
	CreatedAt    string `db:"createdAt"`
	ExpiredAt    string `db:"expiredAt"`
}
