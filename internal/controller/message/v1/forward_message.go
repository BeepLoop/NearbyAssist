package message

import "fmt"

func MessageForwarder() {
	for {
		message := <-broadcastChan

		if socket, ok := clients[message.Receiver]; ok {
			err := socket.WriteJSON(message)
			if err != nil {
				fmt.Printf("error sending message to recipient: %s\n", err.Error())
			}
		} else {
			fmt.Printf("Receiver not found!\n")
			continue
		}

		if socket, ok := clients[message.Sender]; ok {
			err := socket.WriteJSON(message)
			if err != nil {
				fmt.Printf("error sending message to sender: %s\n", err.Error())
			}
		} else {
			fmt.Printf("Sender not found\n")
			continue
		}
	}
}
