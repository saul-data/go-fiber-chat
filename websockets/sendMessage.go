package sockets

import (
	"strconv"
	"time"
)

// ---- Broadcast messages to rooms -----
func SendMessages(room string) {
	cnt := 0
	for {
		select {
		case <-quit:
			// keep messages running even if connection closed
			// return
		default:
			// Do other stuff
			cnt = cnt + 1
			Broadcast <- Message{Room: room, Data: []byte(strconv.Itoa(cnt) + " - I am in room: " + room)}
			time.Sleep(1 * time.Second)
		}
	}
}
