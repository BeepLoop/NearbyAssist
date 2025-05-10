package response

type Vendor struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	ImageUrl  string   `json:"imageUrl"`
	Phone     string   `json:"phone"`
	Rating    string   `json:"rating"`
	Socials   []Social `json:"socials"`
	Expertise []string `json:"expertise"`
	Address   string   `json:"address"`
}
