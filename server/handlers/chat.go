package handlers

//
//import (
//	"bufio"
//	"encoding/json"
//	"log"
//	"time"
//	"wanshow-bingo/sse"
//
//	"github.com/gofiber/fiber/v2"
//	gonanoid "github.com/matoous/go-nanoid/v2"
//)
//
//// ChatMessage represents a single chat message in the system or from a user.
//type ChatMessage struct {
//	ID        string `json:"id"`                 // Unique message ID
//	Type      string `json:"type"`               // "user" | "system"
//	Username  string `json:"username,omitempty"` // Username of the sender, if applicable
//	Message   string `json:"message"`            // Message content
//	Timestamp string `json:"timestamp"`          // Timestamp of when the message was created (HH:MM format)
//}
//
//// ChatStream establishes a Server-Sent Events (SSE) stream for the chat.
//// It registers a new client with the hub, sends keep-alive pings, and
//// broadcasts incoming messages from the hub to the connected client.
////
//// Parameters:
////   - c: Fiber context for the HTTP request/response
////   - hub: The SSE hub that manages all connected clients
////
//// Returns:
////   - error: standard Fiber error if the stream setup fails
//func ChatStream(c *fiber.Ctx, hub *sse.Hub) error {
//
//	// Fiber's BodyStreamWriter provides a bufio.Writer to stream SSE data
//	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
//		defer func() {
//			if r := recover(); r != nil {
//				// Prevent server crash on unexpected panics during streaming
//				log.Println("Recovered in ChatStream", r)
//			}
//		}()
//
//		if w == nil {
//			log.Println("w == nil")
//			return
//		}
//
//		msg := hub.BuildConnectionCount()
//
//		// Initial comment to flush headers and establish the SSE connection
//		_, _ = w.WriteString("data: " + msg + "\n\n")
//		_ = w.Flush()
//
//		// Create a new SSE client (channel for sending messages)
//		client := make(sse.ClientChannel)
//
//		hub.RegisterClient(client)         // Add client to the hub
//		defer hub.UnregisterClient(client) // Remove client when the function exits
//
//		// Keep-alive ticker to prevent client timeouts
//		ticker := time.NewTicker(30 * time.Second)
//		defer ticker.Stop()
//
//		// Main loop: send messages to client or keep-alive pings
//		for {
//			select {
//			case msg, ok := <-client:
//				// New message from the hub
//				if !ok {
//					log.Println("Context Errored")
//					return
//				}
//				_, _ = w.WriteString(msg)
//				_ = w.Flush()
//
//			case <-ticker.C:
//				msg := hub.BuildConnectionCount()
//
//				// Initial comment to flush headers and establish the SSE connection
//				_, _ = w.WriteString("data: " + msg + "\n\n")
//				_ = w.Flush()
//			}
//		}
//	})
//
//	return nil
//}
//
//// SendSystemMessage handles incoming system messages via HTTP POST.
//// It parses the message from the request body, marks it as a system message,
//// timestamps it, and broadcasts it to all connected SSE clients.
////
//// Parameters:
////   - c: Fiber context
////   - hub: The SSE hub
////
//// Returns:
////   - JSON response with status "ok" or error if parsing fails
//func SendSystemMessage(c *fiber.Ctx, hub *sse.Hub) error {
//	var msg ChatMessage
//	if err := c.BodyParser(&msg); err != nil {
//		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
//	}
//
//	msg.Type = "system"
//	msg.Timestamp = time.Now().Format("15:04")
//
//	id, err := gonanoid.New(10)
//
//	if err != nil {
//		return c.Status(500).JSON(fiber.Map{"error": "Error generating unique ID"})
//	}
//
//	msg.ID = id
//
//	payload := sse.SocketEvent{Opcode: "chat.message", Data: msg}
//
//	b, _ := json.Marshal(payload)
//	hub.Broadcast(string(b))
//
//	return c.JSON(fiber.Map{"status": "ok"})
//}
//
//// SendChatMessage handles incoming user messages via HTTP POST.
//// It parses the message, marks it as a user message, timestamps it,
//// and broadcasts it to all connected SSE clients.
////
//// Parameters:
////   - c: Fiber context
////   - hub: The SSE hub
////
//// Returns:
////   - JSON response with status "ok" or error if parsing fails
//func SendChatMessage(c *fiber.Ctx, hub *sse.Hub) error {
//	var msg ChatMessage
//	if err := c.BodyParser(&msg); err != nil {
//		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
//	}
//
//	msg.Type = "user"
//	msg.Timestamp = time.Now().Format("15:04")
//
//	id, err := gonanoid.New(10)
//
//	if err != nil {
//		return c.Status(500).JSON(fiber.Map{"error": "Error generating unique ID"})
//	}
//
//	msg.ID = id
//
//	payload := sse.SocketEvent{Opcode: "chat.message", Data: msg}
//
//	b, _ := json.Marshal(payload)
//	hub.Broadcast(string(b))
//
//	return c.JSON(fiber.Map{"status": "ok"})
//}

//app.Post("/chat/system", func(c *fiber.Ctx) error {
//	password := os.Getenv("HOST_PASSWORD")
//	if c.Get("Authorization") != "Bearer "+password {
//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
//	}
//	return handlers.SendSystemMessage(c, hub)
//})
//
//app.Post("/chat/message", func(c *fiber.Ctx) error {
//	// TODO: Make this endpoint require authorization of some variety.
//	return handlers.SendChatMessage(c, hub)
//})
