package models

type ConversationModel struct {
	UserId            string `json:"userId" db:"id"`
	Name              string `json:"name" db:"name"`
	ImageUrl          string `json:"imageUrl" db:"imageUrl"`
	LastMessage       string `json:"lastMessage" db:"lastMessage"`
	LastMessageSender string `json:"lastMessageSender" db:"lastMessageSender"`
	LastMessageDate   string `json:"lastMessageDate" db:"lastMessageDate"`
	SeenLastMessage   bool   `json:"seen" db:"seenLastMessage"`
}
