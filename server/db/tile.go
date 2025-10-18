package db

import (
	"context"
	"errors"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
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
