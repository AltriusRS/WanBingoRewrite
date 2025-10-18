package chat

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/middleware"
	"wanshow-bingo/sse"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func Post(ctx *fiber.Ctx) error {
	// Get authenticated player
	player, err := middleware.GetPlayerFromContext(ctx)
	if err != nil || player == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if player has chat permission
	if !player.Permissions.HasPermission(models.PermCanSendMessages) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions to send messages",
		})
	}

	// Parse request body
	var msgBody models.MessageRequest
	if err := ctx.BodyParser(&msgBody); err != nil {
		log.Printf("Error parsing message body: %s", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate message content
	if len(msgBody.Contents) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message content cannot be empty",
		})
	}

	if len(msgBody.Contents) > 500 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message too long (max 500 characters)",
		})
	}

	// Moderate content
	moderationResult := utils.ModerateContent(msgBody.Contents)
	if !moderationResult.Allowed {
		log.Printf("Message rejected for user %s: %s", player.ID, moderationResult.Reason)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message contains inappropriate content",
		})
	}

	// Generate message ID
	messageID, err := gonanoid.New(10)
	if err != nil {
		log.Printf("Error generating message ID: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate message ID",
		})
	}

	// Create message object
	message := &models.Message{
		ID:        messageID,
		ShowID:    "", // Will be set by SaveMessage to latest show
		PlayerID:  player.ID,
		Contents:  msgBody.Contents,
		System:    false,
		Replying:  nil,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: nil,
	}

	// Save message to database
	err = db.SaveMessage(context.Background(), message)
	if err != nil {
		log.Printf("Error saving message: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save message",
		})
	}

	// Broadcast message to chat hub
	chatHub := sse.GetChatHub()
	if chatHub != nil {
		chatHub.BroadcastEvent("chat.message", message)
	} else {
		log.Printf("Warning: Chat hub not available for broadcasting")
	}

	// Return success response
	return ctx.JSON(fiber.Map{
		"success":    true,
		"message_id": messageID,
	})
}
