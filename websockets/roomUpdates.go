package sockets

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

var quit = make(chan bool)

// https://github.com/gorilla/websocket/blob/master/examples/chat/client.go

// https://github.com/marcelo-tm/testws/blob/master/main.go

func RoomUpdates(conn *websocket.Conn, room string) {

	// When the function returns, unregister the client and close the connection
	defer func() {
		Unregister <- Subscription{Conn: conn, Room: room}

		// Stop the broadcasting messages
		quit <- true
		conn.Close()
	}()

	// Register the client
	Register <- Subscription{Conn: conn, Room: room}

	// go SecureTimeout()

	for {

		mt, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {

				log.Println("read error:", err)

			}
			return
		}

		log.Println("message received from client:", mt, string(message))

	}

}
