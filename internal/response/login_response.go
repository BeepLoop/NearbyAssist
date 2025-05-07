package response

type DetailedUser struct {
	Id           string      `json:"id"`
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	ImageUrl     string      `json:"imageUrl"`
	IsVerified   bool        `json:"isVerified"`
	IsVendor     bool        `json:"isVendor"`
	Address      string      `json:"address"`
	Phone        string      `json:"phone"`
	Latitude     float64     `json:"latitude"`
	Longitude    float64     `json:"longitude"`
	Expertises   []Expertise `json:"expertises"`
	Socials      []Social    `json:"socials"`
	IsRestricted bool        `json:"isRestricted"`
	DBL          int         `json:"dbl"`
}

type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         DetailedUser `json:"user"`
}

type Social struct {
	Id    string `json:"id"`
	Site  string `json:"site"`
	Title string `json:"title"`
	URL   string `json:"url"`
}
