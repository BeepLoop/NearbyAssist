package response

type Message struct {
	Id        string `json:"id"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Content   string `json:"content"`
	Seen      bool   `json:"seen"`
	SeenAt    string `json:"seenAt"`
	CreatedAt string `json:"createdAt"`
}

type MessageEvent struct {
	Id        string      `json:"id"`
	Sender    MessageUser `json:"sender"`
	Receiver  MessageUser `json:"receiver"`
	Content   string      `json:"content"`
	CreatedAt string      `json:"createdAt"`
	Seen      bool        `json:"seen"`
}

type MessageUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}
