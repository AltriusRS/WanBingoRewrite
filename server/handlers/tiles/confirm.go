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
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type ConfirmTileRequest struct {
	TileID  string `json:"tile_id" validate:"required"`
	Context string `json:"context"`
}

func ConfirmTile(c *fiber.Ctx) error {
	player, err := middleware.GetPlayerFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewApiError("Not authenticated", 401))
	}

	// Check if user is host
	if (player.Permissions & 512) == 0 {
		return c.Status(fiber.StatusForbidden).JSON(utils.NewApiError("Host permission required", 403))
	}

	var req ConfirmTileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Invalid request body", 400))
	}

	if req.TileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Tile ID is required", 400))
	}

	ctx := context.Background()

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get current show", 500))
	}

	id, _ := gonanoid.New(10)

	// Create tile confirmation
	confirmation := &models.TileConfirmation{
		ID:               id,
		ShowID:           latestShow.ID,
		TileID:           req.TileID,
		ConfirmedBy:      &player.ID,
		Context:          nil,
		ConfirmationTime: time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if req.Context != "" {
		confirmation.Context = &req.Context
	}

	// Save to database
	err = db.PersistTileConfirmation(ctx, confirmation)
	if err != nil {
		utils.Debugf("Failed to save confirmation - %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to save confirmation", 500))
	}

	// Get the tile details for the message
	tile, err := db.GetTileByID(ctx, req.TileID)
	if err != nil {
		utils.Debugf("Failed to get tile details - %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get tile details", 500))
	}

	// Create system message
	messageContent := "**TILE CONFIRMED** " + tile.Title
	if req.Context != "" {
		messageContent = "**TILE CONFIRMED**\n\n" + tile.Title + "\n\n> Context: " + req.Context
	}

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
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to send system message", 500))
	}

	// Broadcast message to chat hub
	chatHub := sse.GetChatHub()
	if chatHub != nil {
		chatHub.BroadcastEvent("chat.message", systemMessage)
	} else {
		log.Printf("Warning: Chat hub not available for broadcasting")
	}

	// Broadcast tile confirmation to host hub
	utils.Debugf("[TileConfirm] Broadcasting tile confirmation for tile %s", req.TileID)
	hostHub := sse.GetHostHub()
	if hostHub != nil {
		hostHub.BroadcastEvent("tile.confirm", fiber.Map{
			"tileId": req.TileID,
		})
		utils.Debugf("[TileConfirm] Broadcasted tile.confirm event for tile %s", req.TileID)
	} else {
		utils.Debugf("[TileConfirm] Host hub not available for broadcasting tile confirmation")
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"confirmation": confirmation,
	})
}
