package message

import (
	"fmt"
	query "nearbyassist/internal/db/query/message"
)

func MessageSavior() {
	for {
		message := <-messageChan

		inserted, err := query.NewMessage(message)
		if err != nil {
			fmt.Printf("error saving message: %s\n", err.Error())
			continue
		}

		broadcastChan <- *inserted
	}
}
