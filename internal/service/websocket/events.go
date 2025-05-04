package websocket

import (
	"encoding/json"
	"fmt"
)

type Event string

const (
	EVT_PONG              Event = "pong"
	EVT_MSSG              Event = "message"
	EVT_NOTIF             Event = "notification"
	EVT_SYNC              Event = "sync"
	EVT_BOOKING_COMPLETE  Event = "bookingComplete"
	EVT_BOOKING_CONFIRMED Event = "bookingConfirmed"
	EVT_BOOKING_REJECTED  Event = "bookingRejected"
)

type EventModel struct {
	ReceiverId string      `json:"-"`
	Type       Event       `json:"type"`
	Payload    interface{} `json:"payload"`
}

func (e *EventModel) JsonPrint() {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(b))
}
