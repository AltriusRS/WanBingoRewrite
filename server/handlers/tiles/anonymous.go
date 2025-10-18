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

	// Ensure show has enough tiles
	showTiles, err := db.GetShowTiles(ctx, latestShow.ID)
	if err != nil {
		log.Printf("failed to get show tiles: %v", err)
		return utils.NewApiError("Failed to get show tiles", 0x0502).AsResponse(c)
	}
	if len(showTiles) < db.GetTilesPerShow() {
		err = db.PopulateShowTilesWithRandom(ctx, latestShow.ID)
		if err != nil {
			log.Printf("failed to populate show tiles: %v", err)
			return utils.NewApiError("Failed to populate show tiles", 0x0505).AsResponse(c)
		}
		showTiles, err = db.GetShowTiles(ctx, latestShow.ID)
		if err != nil {
			log.Printf("failed to get show tiles after population: %v", err)
			return utils.NewApiError("Failed to get show tiles", 0x0502).AsResponse(c)
		}
	}

	var availableTileIDs []string
	availableTileIDs = make([]string, len(showTiles))
	for i, showTile := range showTiles {
		availableTileIDs[i] = showTile.TileID
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

	// Calculate potential score
	showTileMap := make(map[string]models.ShowTile)
	for _, st := range showTiles {
		showTileMap[st.TileID] = st
	}
	var potentialScore float64
	for _, tileID := range selectedTiles {
		if st, ok := showTileMap[tileID]; ok {
			potentialScore += st.Score * st.Weight
		}
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
		PotentialScore:         potentialScore,
		RegenerationDiminisher: 1.0,
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
