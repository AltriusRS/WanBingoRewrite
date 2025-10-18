package db

import (
	"context"
	"errors"
	"math/rand"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
	"github.com/matoous/go-nanoid/v2"
)

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
	// Get all available tiles for this show, or fall back to all tiles if no show-specific tiles exist
	showTiles, err := GetShowTiles(ctx, showID, tx...)
	var availableTileIDs []string

	if err != nil || len(showTiles) < 25 {
		// Fall back to all tiles if show-specific tiles are insufficient
		allTiles, err := GetAllTiles(ctx, tx...)
		if err != nil {
			return nil, err
		}
		if len(allTiles) < 25 {
			return nil, errors.New("insufficient tiles available for board generation")
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
		return nil, errors.New("insufficient tiles available for board generation")
	}

	// Randomly select 25 tiles
	selectedTiles := make([]string, 25)
	availableIndices := make([]int, len(availableTileIDs))
	for i := range availableIndices {
		availableIndices[i] = i
	}

	// Fisher-Yates shuffle to select 25 random tiles
	for i := 0; i < 25; i++ {
		j := rand.Intn(len(availableIndices)-i) + i
		availableIndices[i], availableIndices[j] = availableIndices[j], availableIndices[i]
		selectedTiles[i] = availableTileIDs[availableIndices[i]]
	}

	boardID, _ := gonanoid.New(10)

	if len(tx) > 0 {
		_, err = tx[0].Exec(ctx, `
			INSERT INTO boards (id, player_id, show_id, tiles, winner, total_score, potential_score, regeneration_diminisher)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, boardID, playerID, showID, selectedTiles, false, 0, 0, 0)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		_, err = pool.Exec(ctx, `
			INSERT INTO boards (id, player_id, show_id, tiles, winner, total_score, potential_score, regeneration_diminisher)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, boardID, playerID, showID, selectedTiles, false, 0, 0, 0)
	}

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

	// If board doesn't exist, create one
	return CreateBoardForPlayer(ctx, playerID, showID, tx...)
}
