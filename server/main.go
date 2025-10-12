package main

import (
	"encoding/json"
	"log"
	"os"
	"wanshow-bingo/db"
	"wanshow-bingo/handlers"
	"wanshow-bingo/sse"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	// Initialize optional database pool.
	db.Init()

	hub := sse.NewHub()
	go hub.Run()

	// Start WhenPlane aggregate: initial HTTP warmup, then 24/7 websocket, broadcasting to SSE hub
	handlers.StartAggregateFetcher(hub)

	app.Get("/api/tiles", handlers.GetTiles)

	app.Get("/api/aggregate", handlers.GetAggregate)
	app.Post("/api/aggregate", func(c *fiber.Ctx) error {
		password := os.Getenv("HOST_PASSWORD")
		if c.Get("Authorization") != "Bearer "+password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		var payload any

		err := c.BodyParser(&payload)
		if err != nil {
			return err
		}

		event := sse.SocketEvent{Opcode: "whenplane.aggregate", Data: payload}

		b, _ := json.Marshal(event)
		hub.Broadcast(string(b))
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/api/chat/stream", func(c *fiber.Ctx) error {
		return handlers.ChatStream(c, hub)
	})

	app.Post("/api/chat/system", func(c *fiber.Ctx) error {
		password := os.Getenv("HOST_PASSWORD")
		if c.Get("Authorization") != "Bearer "+password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		return handlers.SendSystemMessage(c, hub)
	})

	app.Post("/api/chat/message", func(c *fiber.Ctx) error {
		// TODO: Make this endpoint require authorization of some variety.
		return handlers.SendChatMessage(c, hub)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
