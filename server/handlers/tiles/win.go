package tilerouter

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
	"github.com/google/uuid"
)

func RecordWin(c *fiber.Ctx) error {
	player, err := middleware.GetPlayerFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewApiError("Not authenticated", 401))
	}

	ctx := context.Background()

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get current show", 500))
	}

	// Get player's board
	board, err := db.GetBoardForPlayer(ctx, player.ID, latestShow.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewApiError("Board not found", 404))
	}

	// Check if already won
	if board.Winner {
		return c.Status(fiber.StatusConflict).JSON(utils.NewApiError("Already recorded as winner", 409))
	}

	// Update board to mark as winner
	err = db.UpdateBoardWinner(ctx, board.ID, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to record win", 500))
	}

	// Create system message for the win
	messageContent := "**BINGO WINNER!** " + player.DisplayName + " has won the bingo game!"

	systemMessage := &models.Message{
		ID:        uuid.New().String(),
		ShowID:    latestShow.ID,
		PlayerID:  player.ID,
		Contents:  messageContent,
		System:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save message
	err = db.PersistMessage(ctx, systemMessage)
	if err != nil {
		log.Printf("Failed to send win system message: %v", err)
		// Don't fail the request for this
	}

	// Broadcast message to chat hub
	chatHub := sse.GetChatHub()
	if chatHub != nil {
		chatHub.BroadcastEvent("chat.message", systemMessage)
	} else {
		log.Printf("Warning: Chat hub not available for broadcasting win message")
	}

	// Broadcast win event to host hub
	hostHub := sse.GetHostHub()
	if hostHub != nil {
		hostHub.BroadcastEvent("player.win", fiber.Map{
			"playerId":   player.ID,
			"playerName": player.DisplayName,
			"boardId":    board.ID,
		})
	} else {
		log.Printf("Warning: Host hub not available for broadcasting win event")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Win recorded successfully",
	})
}
