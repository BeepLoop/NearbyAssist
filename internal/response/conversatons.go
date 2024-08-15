package response

type Conversation struct {
	UserId   int    `json:"userId"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}
