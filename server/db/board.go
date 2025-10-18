package db

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
	"github.com/matoous/go-nanoid/v2"
)

var tilesPerShow = 90

// GetTilesPerShow returns the number of tiles per show
func GetTilesPerShow() int {
	return tilesPerShow
}

// GetBoardByPlayerAndShow retrieves a Board for a specific player and show
func GetBoardByPlayerAndShow(ctx context.Context, playerID, showID string, tx ...pgx.Tx) (*models.Board, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, player_id, show_id, tiles, winner, total_score, potential_score, regeneration_diminisher, created_at, updated_at, deleted_at
			FROM boards
			WHERE player_id = $1 AND show_id = $2 AND deleted_at IS NULL
		`, playerID, showID)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, player_id, show_id, tiles, winner, total_score, potential_score, regeneration_diminisher, created_at, updated_at, deleted_at
			FROM boards
			WHERE player_id = $1 AND show_id = $2 AND deleted_at IS NULL
		`, playerID, showID)
	}

	var board models.Board
	err := row.Scan(
		&board.ID, &board.PlayerID, &board.ShowID, &board.Tiles, &board.Winner,
		&board.TotalScore, &board.PotentialScore, &board.RegenerationDiminisher,
		&board.CreatedAt, &board.UpdatedAt, &board.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	return &board, nil
}

// CreateBoardForPlayer creates a new bingo board for a player for the current show
func CreateBoardForPlayer(ctx context.Context, playerID, showID string, tx ...pgx.Tx) (*models.Board, error) {
	// Get the "Show Is Late" tile
	showIsLateTile, err := GetOrCreateShowIsLateTile(ctx, tx...)
	if err != nil {
		return nil, err
	}

	// Ensure "Show Is Late" is in show_tiles
	err = EnsureTileInShowTiles(ctx, showID, showIsLateTile.ID, tx...)
	if err != nil {
		return nil, err
	}

	// Get show-specific tiles
	var showTiles []models.ShowTile
	showTiles, err = GetShowTiles(ctx, showID, tx...)
	if err != nil {
		return nil, errors.New("failed to get show tiles")
	}
	if len(showTiles) < tilesPerShow {
		// Populate with 90 random tiles
		err = PopulateShowTilesWithRandom(ctx, showID, tx...)
		if err != nil {
			return nil, err
		}
		// Get again
		showTiles, err = GetShowTiles(ctx, showID, tx...)
		if err != nil {
			return nil, errors.New("failed to get show tiles after population")
		}
	}
	if len(showTiles) < 25 {
		return nil, errors.New("insufficient show tiles available for board generation")
	}

	availableTileIDs := make([]string, len(showTiles))
	for i, showTile := range showTiles {
		availableTileIDs[i] = showTile.TileID
	}

	// Ensure "Show Is Late" is in show_tiles
	err = EnsureTileInShowTiles(ctx, showID, showIsLateTile.ID, tx...)
	if err != nil {
		return nil, err
	}

	// Get show-specific tiles
	showTiles, err = GetShowTiles(ctx, showID, tx...)
	if err != nil {
		return nil, errors.New("failed to get show tiles")
	}
	if len(showTiles) < tilesPerShow {
		// Populate with tilesPerShow random tiles
		err = PopulateShowTilesWithRandom(ctx, showID, tx...)
		if err != nil {
			return nil, err
		}
		// Get again
		showTiles, err = GetShowTiles(ctx, showID, tx...)
		if err != nil {
			return nil, errors.New("failed to get show tiles after population")
		}
	}
	if len(showTiles) < 25 {
		return nil, errors.New("insufficient show tiles available for board generation")
	}

	availableTileIDs = make([]string, len(showTiles))
	for i, showTile := range showTiles {
		availableTileIDs[i] = showTile.TileID
	}

	if len(availableTileIDs) < 25 {
		return nil, errors.New("insufficient tiles available for board generation")
	}

	// Ensure "Show Is Late" is in the available tiles
	showIsLateIncluded := false
	for _, id := range availableTileIDs {
		if id == showIsLateTile.ID {
			showIsLateIncluded = true
			break
		}
	}
	if !showIsLateIncluded {
		availableTileIDs = append(availableTileIDs, showIsLateTile.ID)
	}

	// Remove "Show Is Late" from available for random selection
	var filteredAvailable []string
	for _, id := range availableTileIDs {
		if id != showIsLateTile.ID {
			filteredAvailable = append(filteredAvailable, id)
		}
	}

	if len(filteredAvailable) < 24 {
		return nil, errors.New("insufficient tiles available for board generation")
	}

	// Randomly select 24 tiles from filtered available
	selectedTiles := make([]string, 25)
	availableIndices := make([]int, len(filteredAvailable))
	for i := range availableIndices {
		availableIndices[i] = i
	}

	// Fisher-Yates shuffle to select 24 random tiles
	for i := 0; i < 24; i++ {
		j := rand.Intn(len(availableIndices)-i) + i
		availableIndices[i], availableIndices[j] = availableIndices[j], availableIndices[i]
		selectedTiles[i] = filteredAvailable[availableIndices[i]]
	}

	// Shift tiles after index 11 to make room for center
	for i := 23; i >= 12; i-- {
		selectedTiles[i+1] = selectedTiles[i]
	}

	// Place "Show Is Late" in the center (index 12)
	selectedTiles[12] = showIsLateTile.ID

	// Calculate potential score
	var potentialScore float64
	if len(showTiles) > 0 {
		showTileMap := make(map[string]models.ShowTile)
		for _, st := range showTiles {
			showTileMap[st.TileID] = st
		}
		for _, tileID := range selectedTiles {
			if st, ok := showTileMap[tileID]; ok {
				potentialScore += st.Score * st.Weight
			}
		}
	} else {
		// Fallback: calculate from tiles table
		for _, tileID := range selectedTiles {
			tile, err := GetTileByID(ctx, tileID)
			if err == nil {
				potentialScore += tile.Score * tile.Weight
			}
		}
	}

	boardID, _ := gonanoid.New(10)

	if len(tx) > 0 {
		_, err = tx[0].Exec(ctx, `
			INSERT INTO boards (id, player_id, show_id, tiles, winner, total_score, potential_score, regeneration_diminisher)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, boardID, playerID, showID, selectedTiles, false, 0, potentialScore, 1.0)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		_, err = pool.Exec(ctx, `
			INSERT INTO boards (id, player_id, show_id, tiles, winner, total_score, potential_score, regeneration_diminisher)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, boardID, playerID, showID, selectedTiles, false, 0, potentialScore, 1.0)
	}
	log.Printf("Created board %s for show %s with potential_score %f, regeneration_diminisher 1.0", boardID, showID, potentialScore)

	if err != nil {
		return nil, err
	}

	// Return the created board
	return GetBoardByPlayerAndShow(ctx, playerID, showID, tx...)
}

// GetOrCreateBoardForPlayer gets a player's board for the current show, creating one if it doesn't exist
func GetOrCreateBoardForPlayer(ctx context.Context, playerID, showID string, tx ...pgx.Tx) (*models.Board, error) {
	board, err := GetBoardByPlayerAndShow(ctx, playerID, showID, tx...)
	if err == nil {
		return board, nil
	}

	log.Println("Failed to get board for player - Generating new board")

	// If board doesn't exist, create one
	return CreateBoardForPlayer(ctx, playerID, showID, tx...)
}

// GetBoardForPlayer gets a player's board for the current show, returning error if not exists
func GetBoardForPlayer(ctx context.Context, playerID, showID string, tx ...pgx.Tx) (*models.Board, error) {
	return GetBoardByPlayerAndShow(ctx, playerID, showID, tx...)
}

// RegenerateBoardForPlayer regenerates a player's board with new tiles and updated diminisher
func RegenerateBoardForPlayer(ctx context.Context, playerID, showID string, newDiminisher float64, tx ...pgx.Tx) (*models.Board, error) {
	// Get the "Show Is Late" tile
	showIsLateTile, err := GetOrCreateShowIsLateTile(ctx, tx...)
	if err != nil {
		return nil, err
	}

	// Ensure "Show Is Late" is in show_tiles
	err = EnsureTileInShowTiles(ctx, showID, showIsLateTile.ID, tx...)
	if err != nil {
		return nil, err
	}

	// Get current board
	board, err := GetBoardByPlayerAndShow(ctx, playerID, showID, tx...)
	if err != nil {
		return nil, err
	}

	// Generate new tiles
	var showTiles []models.ShowTile
	showTiles, err = GetShowTiles(ctx, showID, tx...)
	if err != nil {
		return nil, errors.New("failed to get show tiles")
	}
	if len(showTiles) < tilesPerShow {
		// Populate with tilesPerShow random tiles
		err = PopulateShowTilesWithRandom(ctx, showID, tx...)
		if err != nil {
			return nil, err
		}
		// Get again
		showTiles, err = GetShowTiles(ctx, showID, tx...)
		if err != nil {
			return nil, errors.New("failed to get show tiles after population")
		}
	}
	if len(showTiles) < 25 {
		return nil, errors.New("insufficient show tiles available for board regeneration")
	}

	availableTileIDs := make([]string, len(showTiles))
	for i, showTile := range showTiles {
		availableTileIDs[i] = showTile.TileID
	}

	// Ensure "Show Is Late" is in the available tiles
	showIsLateIncluded := false
	for _, id := range availableTileIDs {
		if id == showIsLateTile.ID {
			showIsLateIncluded = true
			break
		}
	}
	if !showIsLateIncluded {
		availableTileIDs = append(availableTileIDs, showIsLateTile.ID)
	}

	// Remove "Show Is Late" from available for random selection
	var filteredAvailable []string
	for _, id := range availableTileIDs {
		if id != showIsLateTile.ID {
			filteredAvailable = append(filteredAvailable, id)
		}
	}

	if len(filteredAvailable) < 24 {
		return nil, errors.New("insufficient tiles available for board regeneration")
	}

	// Randomly select 24 tiles from filtered available
	selectedTiles := make([]string, 25)
	availableIndices := make([]int, len(filteredAvailable))
	for i := range availableIndices {
		availableIndices[i] = i
	}

	for i := 0; i < 24; i++ {
		j := rand.Intn(len(availableIndices)-i) + i
		availableIndices[i], availableIndices[j] = availableIndices[j], availableIndices[i]
		selectedTiles[i] = filteredAvailable[availableIndices[i]]
	}

	// Shift tiles after index 11 to make room for center
	for i := 23; i >= 12; i-- {
		selectedTiles[i+1] = selectedTiles[i]
	}

	// Place "Show Is Late" in the center (index 12)
	selectedTiles[12] = showIsLateTile.ID

	// Calculate potential score
	var potentialScore float64
	if len(showTiles) > 0 {
		showTileMap := make(map[string]models.ShowTile)
		for _, st := range showTiles {
			showTileMap[st.TileID] = st
		}
		for _, tileID := range selectedTiles {
			if st, ok := showTileMap[tileID]; ok {
				potentialScore += st.Score * st.Weight
			}
		}
	} else {
		// Fallback: calculate from tiles table
		for _, tileID := range selectedTiles {
			tile, err := GetTileByID(ctx, tileID)
			if err == nil {
				potentialScore += tile.Score * tile.Weight
			}
		}
	}

	// Update board
	if len(tx) > 0 {
		_, err = tx[0].Exec(ctx, `
			UPDATE boards
			SET tiles = $1, winner = false, total_score = 0, potential_score = $2, regeneration_diminisher = $3, updated_at = NOW()
			WHERE id = $4
		`, selectedTiles, potentialScore, newDiminisher, board.ID)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		_, err = pool.Exec(ctx, `
			UPDATE boards
			SET tiles = $1, winner = false, total_score = 0, potential_score = $2, regeneration_diminisher = $3, updated_at = NOW()
			WHERE id = $4
		`, selectedTiles, potentialScore, newDiminisher, board.ID)
	}

	if err != nil {
		return nil, err
	}

	// Return updated board
	return GetBoardByPlayerAndShow(ctx, playerID, showID, tx...)
}
