package ws_quicktasking

import (
	"fmt"
)

func Initialize() {
	fmt.Println("Connecting to quicktask websocket...")
	success := make(chan bool)
	go handleWebsocket(success)
	didSucceed := <-success
	if didSucceed {
		fmt.Println("Successfully authenticated to quicktask websocket.")
	}
}
