package sockets

import (
	"log"

	"github.com/gofiber/websocket/v2"
	cmap "github.com/orcaman/concurrent-map"
)

type client struct{} // Add more data to this type if needed

var clients = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
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

	RoomswithClients := cmap.New()

	for {
		select {
		case connection := <-Register:
			clients[connection.Conn] = client{}

			var clientsperRoom = []*websocket.Conn{}

			// Retrieve list of clients in the room.
			if tmp, ok := RoomswithClients.Get(connection.Room); ok {
				clientsperRoom = tmp.([]*websocket.Conn)
				clientsperRoom = append(clientsperRoom, connection.Conn)
				RoomswithClients.Set(connection.Room, clientsperRoom)
			} else {
				clientsperRoom = append(clientsperRoom, connection.Conn)
				RoomswithClients.Set(connection.Room, clientsperRoom)
			}

			// Show clients belonging to which room
			log.Println("connections registered:")

			// if tmp, ok := RoomswithClients.Get(connection.Room); ok {
			// 	clientsperRoom = tmp.([]*websocket.Conn)
			// }
			for i, rooms := range RoomswithClients {
				log.Println("Reg: ", i, rooms)
			}

		case message := <-Broadcast:

			// log.Println("message received:", string(message.Data), message.Room)

			var clientsperRoom = []*websocket.Conn{}

			// Send the message to all clients
			if tmp, ok := RoomswithClients.Get(message.Room); ok {
				clientsperRoom = tmp.([]*websocket.Conn)
			}

			for _, connection := range clientsperRoom {

				// If an error occurs close the client
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message.Data)); err != nil {
					log.Println("write error:", err)

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(clients, connection)
				}
			}

		case connection := <-Unregister:
			// Remove the client from the hub
			delete(clients, connection.Conn)

			log.Println("connection unregistered")
		}
	}
}
