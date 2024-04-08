package message

import (
	"fmt"
)

func MessageSavior() {
	for {
		message := <-messageChan

		err := message.Save()
		if err != nil {
			fmt.Printf("error saving message: %s\n", err.Error())
			continue
		}

		broadcastChan <- message
	}
}
