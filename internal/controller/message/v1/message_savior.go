package message

import (
	"fmt"
	"nearbyassist/internal/db/query/message"
)

func MessageSavior() {
	for {
		message := <-messageChan

		inserted, err := message_query.NewMessage(message)
		if err != nil {
			fmt.Printf("error saving message: %s\n", err.Error())
			continue
		}

		broadcastChan <- *inserted
	}
}
