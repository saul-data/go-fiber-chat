package main

import (
	sockets "saul-data/chat/websockets"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use(logger.New(
		logger.Config{
			Format: "✨ Latency: ${latency} Time:${time} Status: ${status} Path:${path} \n",
		}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello 👋!")
	})

	app.Get("/api/:room", func(c *fiber.Ctx) error {

		room := string(c.Params("room"))
		sockets.Broadcast <- sockets.Message{Room: room, Data: []byte("API: I am in room: " + room)}
		return c.SendString("Hello 👋! I am in room " + room)
	})

	// Websockets
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Run the sockets hub
	go sockets.RunHub()

	// Broadcast messages to subscribers in channels room1 and room2
	go sockets.SendMessages("room1")
	go sockets.SendMessages("room2")

	app.Get("/ws/rooms/:room", websocket.New(func(c *websocket.Conn) {

		room := string(c.Params("room"))
		sockets.RoomUpdates(c, room)
	}))

	app.Listen(":9000")
}
