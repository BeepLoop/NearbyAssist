package response

type DetailedUser struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	ImageUrl   string      `json:"imageUrl"`
	IsVerified bool        `json:"isVerified"`
	IsVendor   bool        `json:"isVendor"`
	Address    string      `json:"address"`
	Latitude   float64     `json:"latitude"`
	Longitude  float64     `json:"longitude"`
	Expertises []Expertise `json:"expertises"`
}

type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         DetailedUser `json:"user"`
}
