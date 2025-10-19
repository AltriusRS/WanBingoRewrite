package host

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

type CreateTileRequest struct {
	Text     string                 `json:"text" validate:"required"`
	Category *string                `json:"category"`
	Weight   float64                `json:"weight"`
	Score    float64                `json:"score"`
	Settings map[string]interface{} `json:"settings"`
}

type UpdateTileRequest struct {
	Text     *string                `json:"text"`
	Category *string                `json:"category"`
	Weight   *float64               `json:"weight"`
	Score    *float64               `json:"score"`
	Settings map[string]interface{} `json:"settings"`
}

func init() {
	utils.RegisterRouter("/host", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	// Require host authentication
	host := router.Group("", middleware.AuthMiddleware, requireHost)

	host.Get("/tile-stats", GetTileStats)
	host.Get("/tiles", GetTiles)
	host.Get("/confirmed-tiles", GetConfirmedTiles)
	host.Post("/tile-locks", LockTile)
	host.Post("/tile-unlocks", UnlockTile)
	host.Delete("/confirmed-tiles/:tileId", RevokeConfirmation)
	host.Post("/tiles", CreateTile)
	host.Patch("/tiles/:id", UpdateTile)
	host.Delete("/tiles/:id", DeleteTile)
}

func requireHost(c *fiber.Ctx) error {
	player, err := middleware.GetPlayerFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewApiError("Not authenticated", 401))
	}

	if (player.Permissions & 512) == 0 {
		return c.Status(fiber.StatusForbidden).JSON(utils.NewApiError("Host permission required", 403))
	}

	return c.Next()
}

func GetTiles(c *fiber.Ctx) error {
	ctx := context.Background()
	tiles, err := db.GetAllTiles(ctx)
	if err != nil {
		log.Printf("Failed to get tiles: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get tiles", 500))
	}

	return c.JSON(tiles)
}

func CreateTile(c *fiber.Ctx) error {
	var req CreateTileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Invalid request body", 400))
	}

	if req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Tile text is required", 400))
	}

	ctx := context.Background()

	tileID, _ := gonanoid.New(10)
	tile := &models.Tile{
		ID:        tileID,
		Title:     req.Text,
		Category:  req.Category,
		Weight:    req.Weight,
		Score:     req.Score,
		Settings:  req.Settings,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := db.PersistTile(ctx, tile)
	if err != nil {
		log.Printf("Failed to create tile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to create tile", 500))
	}

	return c.Status(fiber.StatusCreated).JSON(tile)
}

func UpdateTile(c *fiber.Ctx) error {
	tileID := c.Params("id")
	if tileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Tile ID is required", 400))
	}

	var req UpdateTileRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Invalid request body", 400))
	}

	log.Printf("UpdateTile request: %+v", req)

	ctx := context.Background()

	// Get existing tile
	tile, err := db.GetTileByID(ctx, tileID)
	if err != nil {
		log.Printf("Failed to get tile %s: %v", tileID, err)
		return c.Status(fiber.StatusNotFound).JSON(utils.NewApiError("Tile not found", 404))
	}

	log.Printf("Existing tile settings: %+v", tile.Settings)

	// Update fields
	if req.Text != nil {
		tile.Title = *req.Text
	}
	if req.Category != nil {
		tile.Category = req.Category
	}
	if req.Weight != nil {
		tile.Weight = *req.Weight
	}
	if req.Score != nil {
		tile.Score = *req.Score
	}
	if req.Settings != nil {
		log.Printf("Updating settings to: %+v", req.Settings)
		tile.Settings = req.Settings
	}

	tile.UpdatedAt = time.Now()

	log.Printf("Tile before persist: %+v", tile)

	err = db.PersistTile(ctx, tile)
	if err != nil {
		log.Printf("Failed to update tile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to update tile", 500))
	}

	log.Printf("Tile updated successfully")

	// Fetch the updated tile to verify
	updatedTile, err := db.GetTileByID(ctx, tileID)
	if err != nil {
		log.Printf("Failed to fetch updated tile: %v", err)
	} else {
		log.Printf("Updated tile settings: %+v", updatedTile.Settings)
	}

	return c.JSON(tile)
}

func DeleteTile(c *fiber.Ctx) error {
	tileID := c.Params("id")
	if tileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Tile ID is required", 400))
	}

	ctx := context.Background()

	err := db.DeleteTile(ctx, tileID)
	if err != nil {
		log.Printf("Failed to delete tile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to delete tile", 500))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

func GetConfirmedTiles(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		log.Printf("Failed to get latest show: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get latest show", 500))
	}

	// Get confirmed tiles for this show
	confirmed, err := db.GetTileConfirmationsForShow(ctx, latestShow.ID)
	if err != nil {
		log.Printf("Failed to get confirmed tiles: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get confirmed tiles", 500))
	}

	// Return just the tile IDs
	tileIDs := make([]string, len(confirmed))
	for i, conf := range confirmed {
		tileIDs[i] = conf.TileID
	}

	return c.JSON(tileIDs)
}

type LockRequest struct {
	TileID string `json:"tile_id" validate:"required"`
}

func LockTile(c *fiber.Ctx) error {
	var req LockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Invalid request body", 400))
	}

	player, err := middleware.GetPlayerFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewApiError("Not authenticated", 401))
	}

	utils.Debugf("[HostTiles] LockTile: locking tile %s for user %s", req.TileID, player.DisplayName)

	hostHub := sse.GetHostHub()
	if hostHub != nil {
		hostHub.BroadcastEvent("tile.lock", fiber.Map{
			"tileId": req.TileID,
			"user":   player.DisplayName,
		})
		utils.Debugf("[HostTiles] LockTile: broadcasted lock event for tile %s", req.TileID)
	} else {
		utils.Debugf("[HostTiles] LockTile: host hub not available")
	}

	return c.JSON(fiber.Map{"success": true})
}

func UnlockTile(c *fiber.Ctx) error {
	var req LockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Invalid request body", 400))
	}

	utils.Debugf("[HostTiles] UnlockTile: unlocking tile %s", req.TileID)

	hostHub := sse.GetHostHub()
	if hostHub != nil {
		hostHub.BroadcastEvent("tile.unlock", fiber.Map{
			"tileId": req.TileID,
		})
		utils.Debugf("[HostTiles] UnlockTile: broadcasted unlock event for tile %s", req.TileID)
	} else {
		utils.Debugf("[HostTiles] UnlockTile: host hub not available")
	}

	return c.JSON(fiber.Map{"success": true})
}

func RevokeConfirmation(c *fiber.Ctx) error {
	tileID := c.Params("tileId")
	if tileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Tile ID is required", 400))
	}

	ctx := context.Background()

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		log.Printf("Failed to get latest show: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get latest show", 500))
	}

	// Delete the confirmation
	err = db.DeleteTileConfirmation(ctx, latestShow.ID, tileID)
	if err != nil {
		log.Printf("Failed to revoke confirmation: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to revoke confirmation", 500))
	}

	// Broadcast revoke
	utils.Debugf("[HostTiles] RevokeConfirmation: revoking tile %s", tileID)
	hostHub := sse.GetHostHub()
	if hostHub != nil {
		hostHub.BroadcastEvent("tile.revoke", fiber.Map{
			"tileId": tileID,
		})
		utils.Debugf("[HostTiles] RevokeConfirmation: broadcasted revoke event for tile %s", tileID)
	} else {
		utils.Debugf("[HostTiles] RevokeConfirmation: host hub not available")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

func GetTileStats(c *fiber.Ctx) error {
	// This would require implementing stats calculation
	// For now, return empty array
	return c.JSON([]interface{}{})
}

func PostTestMessage(c *fiber.Ctx) error {
	// Get authenticated player
	player, err := middleware.GetPlayerFromContext(c)
	if err != nil || player == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if player is host (has host permission)
	if !player.Permissions.HasPermission(models.PermCanHost) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Host permission required",
		})
	}

	// Parse request body
	var req struct {
		Message string `json:"message"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message cannot be empty",
		})
	}

	// Generate message ID
	messageID, err := gonanoid.New(10)
	if err != nil {
		log.Printf("Error generating message ID: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate message ID",
		})
	}

	// Create system message
	message := &models.Message{
		ID:        messageID,
		ShowID:    "", // Will be set by SaveMessage to latest show
		PlayerID:  "SYSTEM",
		Contents:  req.Message,
		System:    true,
		Replying:  nil,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: nil,
	}

	// Save message to database
	err = db.SaveMessage(context.Background(), message)
	if err != nil {
		log.Printf("Error saving test message: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save message",
		})
	}

	// Broadcast message to chat hub
	chatHub := sse.GetChatHub()
	if chatHub != nil {
		chatHub.BroadcastEvent("chat.message", message)
	} else {
		log.Printf("Warning: Chat hub not available for broadcasting test message")
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"message_id": messageID,
	})
}
