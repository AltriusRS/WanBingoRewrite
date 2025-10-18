package db

import (
	"context"
	"errors"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
)

// GetShowTiles retrieves all tiles for a specific show
func GetShowTiles(ctx context.Context, showID string, tx ...pgx.Tx) ([]models.ShowTile, error) {
	var rows pgx.Rows
	var err error

	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, `
			SELECT show_id, tile_id, weight, score, created_at, updated_at, deleted_at
			FROM show_tiles
			WHERE show_id = $1 AND deleted_at IS NULL
			ORDER BY created_at
		`, showID)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		rows, err = pool.Query(ctx, `
			SELECT show_id, tile_id, weight, score, created_at, updated_at, deleted_at
			FROM show_tiles
			WHERE show_id = $1 AND deleted_at IS NULL
			ORDER BY created_at
		`, showID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var showTiles []models.ShowTile
	for rows.Next() {
		var showTile models.ShowTile
		err := rows.Scan(
			&showTile.ShowID, &showTile.TileID, &showTile.Weight, &showTile.Score,
			&showTile.CreatedAt, &showTile.UpdatedAt, &showTile.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		showTiles = append(showTiles, showTile)
	}

	return showTiles, rows.Err()
}

// GetShowTileIDs returns just the tile IDs for a show (convenience function)
func GetShowTileIDs(ctx context.Context, showID string, tx ...pgx.Tx) ([]string, error) {
	showTiles, err := GetShowTiles(ctx, showID, tx...)
	if err != nil {
		return nil, err
	}

	tileIDs := make([]string, len(showTiles))
	for i, showTile := range showTiles {
		tileIDs[i] = showTile.TileID
	}

	return tileIDs, nil
}
