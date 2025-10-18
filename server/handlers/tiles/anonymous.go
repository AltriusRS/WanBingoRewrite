package tilerouter

import (
	"context"
	"log"
	"math/rand"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// GetAnonymousBoard generates a bingo board without persisting it to the database
func GetAnonymousBoard(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		log.Printf("failed to get latest show: %v", err)
		return utils.NewApiError("Failed to get latest show", 0x0501).AsResponse(c)
	}

	// Get all available tiles for this show, or fall back to all tiles
	showTiles, err := db.GetShowTiles(ctx, latestShow.ID)
	var availableTileIDs []string

	if err != nil || len(showTiles) < 25 {
		// Fall back to all tiles if show-specific tiles are insufficient
		allTiles, err := db.GetAllTiles(ctx)
		if err != nil {
			log.Printf("failed to get all tiles: %v", err)
			return utils.NewApiError("Failed to get tiles", 0x0502).AsResponse(c)
		}
		if len(allTiles) < 25 {
			return utils.NewApiError("Insufficient tiles available for board generation", 0x0503).AsResponse(c)
		}
		availableTileIDs = make([]string, len(allTiles))
		for i, tile := range allTiles {
			availableTileIDs[i] = tile.ID
		}
	} else {
		availableTileIDs = make([]string, len(showTiles))
		for i, showTile := range showTiles {
			availableTileIDs[i] = showTile.TileID
		}
	}

	if len(availableTileIDs) < 25 {
		return utils.NewApiError("Insufficient tiles available for board generation", 0x0504).AsResponse(c)
	}

	// Randomly select 25 tiles using Fisher-Yates shuffle
	selectedTiles := make([]string, 25)
	availableIndices := make([]int, len(availableTileIDs))
	for i := range availableIndices {
		availableIndices[i] = i
	}

	for i := 0; i < 25; i++ {
		j := rand.Intn(len(availableIndices)-i) + i
		availableIndices[i], availableIndices[j] = availableIndices[j], availableIndices[i]
		selectedTiles[i] = availableTileIDs[availableIndices[i]]
	}

	// Get tile details for the selected tiles
	tileDetails := make([]map[string]interface{}, len(selectedTiles))
	for i, tileID := range selectedTiles {
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

	// Create a temporary board object (not persisted)
	board := models.Board{
		ID:                     "anonymous-" + time.Now().Format("20060102150405"),
		PlayerID:               "anonymous",
		ShowID:                 latestShow.ID,
		Tiles:                  selectedTiles,
		Winner:                 false,
		TotalScore:             0,
		PotentialScore:         0,
		RegenerationDiminisher: 0,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
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
		"is_anonymous":    true,
	})
}
