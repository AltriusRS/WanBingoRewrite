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
		"board_id":                board.ID,
		"show_id":                 board.ShowID,
		"player_id":               board.PlayerID,
		"tiles":                   tileDetails,
		"winner":                  board.Winner,
		"total_score":             board.TotalScore,
		"potential_score":         board.PotentialScore,
		"regeneration_diminisher": board.RegenerationDiminisher,
		"created_at":              board.CreatedAt,
	})
}

// RegenerateMyBoard regenerates the player's board with a penalty
func RegenerateMyBoard(c *fiber.Ctx) error {
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

	// Get current board
	board, err := db.GetBoardForPlayer(ctx, player.ID, latestShow.ID)
	if err != nil {
		log.Printf("failed to get board for player %s: %v", player.ID, err)
		return utils.NewApiError("Failed to get board", 0x0404).AsResponse(c)
	}

	// Check regeneration limit
	if board.RegenerationDiminisher <= 0.7 {
		return utils.NewApiError("You can only regenerate your board three times", 0x0405).AsResponse(c)
	}

	// Calculate new diminisher
	var newDiminisher float64
	if board.RegenerationDiminisher == 1 {
		newDiminisher = 0.9
	} else if board.RegenerationDiminisher == 0.9 {
		newDiminisher = 0.8
	} else if board.RegenerationDiminisher == 0.8 {
		newDiminisher = 0.7
	} else {
		return utils.NewApiError("Invalid regeneration state", 0x0406).AsResponse(c)
	}

	// Regenerate board
	newBoard, err := db.RegenerateBoardForPlayer(ctx, player.ID, latestShow.ID, newDiminisher)
	if err != nil {
		log.Printf("failed to regenerate board for player %s: %v", player.ID, err)
		return utils.NewApiError("Failed to regenerate board", 0x0407).AsResponse(c)
	}

	// Get tile details for the new board
	tileDetails := make([]map[string]interface{}, len(newBoard.Tiles))
	for i, tileID := range newBoard.Tiles {
		tile, err := db.GetTileByID(ctx, tileID)
		if err != nil {
			log.Printf("failed to get tile %s: %v", tileID, err)
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
		"board_id":                newBoard.ID,
		"show_id":                 newBoard.ShowID,
		"player_id":               newBoard.PlayerID,
		"tiles":                   tileDetails,
		"winner":                  newBoard.Winner,
		"total_score":             newBoard.TotalScore,
		"potential_score":         newBoard.PotentialScore,
		"regeneration_diminisher": newBoard.RegenerationDiminisher,
		"created_at":              newBoard.CreatedAt,
	})
}
