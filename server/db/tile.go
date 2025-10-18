package db

import (
	"context"
	"errors"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
	"github.com/matoous/go-nanoid/v2"
)

// GetAllTiles retrieves all tiles from the database
func GetAllTiles(ctx context.Context, tx ...pgx.Tx) ([]models.Tile, error) {
	var rows pgx.Rows
	var err error

	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, `
			SELECT id, title, category, last_drawn, weight, score, created_by, settings, created_at, updated_at, deleted_at
			FROM tiles
			WHERE deleted_at IS NULL
			ORDER BY created_at
		`)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		rows, err = pool.Query(ctx, `
			SELECT id, title, category, last_drawn, weight, score, created_by, settings, created_at, updated_at, deleted_at
			FROM tiles
			WHERE deleted_at IS NULL
			ORDER BY created_at
		`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tiles []models.Tile
	for rows.Next() {
		var tile models.Tile
		err := rows.Scan(
			&tile.ID, &tile.Title, &tile.Category, &tile.LastDrawn, &tile.Weight, &tile.Score,
			&tile.CreatedBy, &tile.Settings, &tile.CreatedAt, &tile.UpdatedAt, &tile.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		tiles = append(tiles, tile)
	}

	return tiles, rows.Err()
}

// GetTileByID retrieves a single tile by ID
func GetTileByID(ctx context.Context, id string, tx ...pgx.Tx) (*models.Tile, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, title, category, last_drawn, weight, score, created_by, settings, created_at, updated_at, deleted_at
			FROM tiles
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, title, category, last_drawn, weight, score, created_by, settings, created_at, updated_at, deleted_at
			FROM tiles
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	}

	var tile models.Tile
	err := row.Scan(
		&tile.ID, &tile.Title, &tile.Category, &tile.LastDrawn, &tile.Weight, &tile.Score,
		&tile.CreatedBy, &tile.Settings, &tile.CreatedAt, &tile.UpdatedAt, &tile.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("tile not found")
		}
		return nil, err
	}

	return &tile, nil
}

// GetOrCreateShowIsLateTile ensures the "Show Is Late" tile exists and returns it
func GetOrCreateShowIsLateTile(ctx context.Context, tx ...pgx.Tx) (*models.Tile, error) {
	// First, try to find it by title
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, title, category, last_drawn, weight, score, created_by, settings, created_at, updated_at, deleted_at
			FROM tiles
			WHERE title = $1 AND deleted_at IS NULL
		`, "Show Is Late")
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, title, category, last_drawn, weight, score, created_by, settings, created_at, updated_at, deleted_at
			FROM tiles
			WHERE title = $1 AND deleted_at IS NULL
		`, "Show Is Late")
	}

	var tile models.Tile
	err := row.Scan(
		&tile.ID, &tile.Title, &tile.Category, &tile.LastDrawn, &tile.Weight, &tile.Score,
		&tile.CreatedBy, &tile.Settings, &tile.CreatedAt, &tile.UpdatedAt, &tile.DeletedAt,
	)

	if err == nil {
		return &tile, nil
	}

	if err != pgx.ErrNoRows {
		return nil, err
	}

	// Tile doesn't exist, create it
	tileID, _ := gonanoid.New(10)
	category := "Late"
	weight := 1.0
	score := 5.0
	settings := map[string]interface{}{}

	if len(tx) > 0 {
		_, err = tx[0].Exec(ctx, `
			INSERT INTO tiles (id, title, category, weight, score, settings)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, tileID, "Show Is Late", category, weight, score, settings)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		_, err = pool.Exec(ctx, `
			INSERT INTO tiles (id, title, category, weight, score, settings)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, tileID, "Show Is Late", category, weight, score, settings)
	}

	if err != nil {
		return nil, err
	}

	// Return the created tile
	return GetTileByID(ctx, tileID, tx...)
}

// PersistTileConfirmation saves a tile confirmation to the database
func PersistTileConfirmation(ctx context.Context, confirmation *models.TileConfirmation, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, `
			INSERT INTO tile_confirmations (id, show_id, tile_id, confirmed_by, context, confirmation_time, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, confirmation.ID, confirmation.ShowID, confirmation.TileID, confirmation.ConfirmedBy, confirmation.Context, confirmation.ConfirmationTime, confirmation.CreatedAt, confirmation.UpdatedAt)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, `
			INSERT INTO tile_confirmations (id, show_id, tile_id, confirmed_by, context, confirmation_time, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, confirmation.ID, confirmation.ShowID, confirmation.TileID, confirmation.ConfirmedBy, confirmation.Context, confirmation.ConfirmationTime, confirmation.CreatedAt, confirmation.UpdatedAt)
		return err
	}
}
