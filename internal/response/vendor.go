package response

type Vendor struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	ImageUrl  string   `json:"imageUrl"`
	Phone     string   `json:"phone"`
	Rating    string   `json:"rating"`
	Socials   []string `json:"socials"`
	Expertise []string `json:"expertise"`
}
