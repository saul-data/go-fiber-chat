package sockets

import (
	"log"

	"github.com/gofiber/websocket/v2"
	cmap "github.com/orcaman/concurrent-map"
)

type client struct{} // Add more data to this type if needed

var Register = make(chan Subscription)
var Broadcast = make(chan Message)
var Unregister = make(chan Subscription)

type Message struct {
	Data []byte
	Room string
}

type Subscription struct {
	Conn *websocket.Conn
	Room string
}

func RunHub() {

	// Rooms with list of clients
	RoomswithClients := cmap.New()

	for {
		select {
		case connection := <-Register:

			var clients = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast

			// Retrieve list of clients in the room.
			if tmp, ok := RoomswithClients.Get(connection.Room); ok {

				// Existing room add connection
				clients = tmp.(map[*websocket.Conn]client)
				clients[connection.Conn] = client{}

				RoomswithClients.Set(connection.Room, clients)
			} else {

				// New room add connection
				clients[connection.Conn] = client{}
				RoomswithClients.Set(connection.Room, clients)
			}

			// Show clients belonging to which room
			log.Println("connections registered:")

			// if tmp, ok := RoomswithClients.Get(connection.Room); ok {
			// 	clientsperRoom = tmp.([]*websocket.Conn)
			// }
			// for i, rooms := range RoomswithClients {
			// 	log.Println("Reg: ", i, rooms)
			// }

		case message := <-Broadcast:

			// log.Println("message received:", string(message.Data), message.Room)
			var clients = make(map[*websocket.Conn]client)

			// Send the message to all clients in the room provided
			if tmp, ok := RoomswithClients.Get(message.Room); ok {
				clients = tmp.(map[*websocket.Conn]client)

				for _, connection := range clients {
					log.Println("Connections:", connection)
				}

				for connection := range clients {

					// If an error occurs close the client

					if err := connection.WriteMessage(websocket.TextMessage, []byte(message.Data)); err != nil {
						log.Println("write error:", err)

						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
						delete(clients, connection)
						RoomswithClients.Set(message.Room, clients)
					}
				}
			}

		case connection := <-Unregister:
			// Remove the client from the hub

			if tmp, ok := RoomswithClients.Get(connection.Room); ok {
				var clients = make(map[*websocket.Conn]client)
				clients = tmp.(map[*websocket.Conn]client)
				delete(clients, connection.Conn)
				RoomswithClients.Set(connection.Room, clients)
			}

			log.Println("connection unregistered")
		}
	}
}
