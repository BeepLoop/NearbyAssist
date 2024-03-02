package message

import "fmt"

func MessageForwarder() {
	for {
		message := <-broadcastChan
		fmt.Printf("message: %v\n", message)

		if socket, ok := clients[message.Reciever]; ok {
			err := socket.WriteJSON(message)
			if err != nil {
				fmt.Printf("error sending message to recipient: %s\n", err.Error())
			}
		}
		fmt.Printf("sent to reciever\n")

		if socket, ok := clients[message.Sender]; ok {
			err := socket.WriteJSON(message)
			if err != nil {
				fmt.Printf("error sending message to sender: %s\n", err.Error())
			}
		}
		fmt.Printf("sent to sender\n")
	}
}
