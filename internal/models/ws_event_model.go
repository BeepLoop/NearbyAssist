package models

type EventType string

const (
	EVT_MSSG  = "message"
	EVT_NOTIF = "notification"
	EVT_SYNC  = "sync"
)

type WsEventModel struct {
	Type    EventType
	Payload interface{} `json:"payload"`
}
