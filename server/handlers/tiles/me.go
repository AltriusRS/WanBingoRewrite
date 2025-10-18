package tilerouter

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// GetMyBoard returns the current player's bingo board, creating one if it doesn't exist
func GetMyBoard(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get authenticated player from context
	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0401).AsResponse(c)
	}

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		log.Printf("failed to get latest show: %v", err)
		return utils.NewApiError("Failed to get latest show", 0x0402).AsResponse(c)
	}

	// Get or create board for this player and show
	board, err := db.GetOrCreateBoardForPlayer(ctx, player.ID, latestShow.ID)
	if err != nil {
		log.Printf("failed to get/create board for player %s: %v", player.ID, err)
		return utils.NewApiError("Failed to get/create board", 0x0403).AsResponse(c)
	}

	// Get tile details for the board
	tileDetails := make([]map[string]interface{}, len(board.Tiles))
	for i, tileID := range board.Tiles {
		tile, err := db.GetTileByID(ctx, tileID)
		if err != nil {
			log.Printf("failed to get tile %s: %v", tileID, err)
			// Continue with partial data rather than failing completely
			tileDetails[i] = map[string]interface{}{
				"id":    tileID,
				"error": "tile not found",
			}
			continue
		}
		tileDetails[i] = map[string]interface{}{
			"id":       tile.ID,
			"title":    tile.Title,
			"category": tile.Category,
			"weight":   tile.Weight,
			"score":    tile.Score,
		}
	}

	return c.JSON(fiber.Map{
		"board_id":        board.ID,
		"show_id":         board.ShowID,
		"player_id":       board.PlayerID,
		"tiles":           tileDetails,
		"winner":          board.Winner,
		"total_score":     board.TotalScore,
		"potential_score": board.PotentialScore,
		"created_at":      board.CreatedAt,
	})
}
