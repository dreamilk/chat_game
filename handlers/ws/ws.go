package ws

import "fmt"

func Receive(userID string, msg []byte) {
	fmt.Println(userID, string(msg))
}
