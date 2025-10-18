package chat

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/sse"

	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func PostSystem(ctx *fiber.Ctx) error {
	var msgBody models.MessageRequest

	id, err := gonanoid.New(10)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Error generating unique ID"})
	}

	if err := ctx.BodyParser(&msgBody); err != nil {
		log.Printf("Error parsing body: %s", err)
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	msg := models.Message{
		ID:        id,
		ShowID:    "",
		PlayerID:  "SYSTEM",
		Contents:  msgBody.Contents,
		System:    true,
		Replying:  nil,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: nil,
	}

	err = db.SaveMessage(context.Background(), &msg)

	if err != nil {
		log.Printf("Error posting message: %s", err)
		return ctx.Status(500).JSON(fiber.Map{"error": "error posting message"})
	}

	ch := sse.GetChatHub()

	if ch == nil {
		log.Printf("Error parsing body: ch is nil")
		return ctx.Status(500).JSON(fiber.Map{"error": "error posting message"})
	}

	ch.BroadcastEvent("chat.message", msg)

	return nil
}
