package models

type ConversationModel struct {
	UserId          string `json:"userId" db:"id"`
	Name            string `json:"name" db:"name"`
	ImagaUrl        string `json:"imageUrl" db:"imageUrl"`
	LastMessage     string `json:"lastMessage" db:"lastMessage"`
	LastMessageDate string `json:"lastMessageDate" db:"lastMessageDate"`
}
