package db

import (
	"context"
	"errors"
	"log"
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

// EnsureTileInShowTiles ensures a tile is associated with a show
func EnsureTileInShowTiles(ctx context.Context, showID, tileID string, tx ...pgx.Tx) error {
	// Check if already exists
	var exists bool
	if len(tx) > 0 {
		err := tx[0].QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM show_tiles WHERE show_id = $1 AND tile_id = $2)
		`, showID, tileID).Scan(&exists)
		if err != nil {
			return err
		}
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		err := pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM show_tiles WHERE show_id = $1 AND tile_id = $2)
		`, showID, tileID).Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists {
		return nil
	}

	log.Printf("Creating new show tile: %s - %s", showID, tileID)

	// Insert
	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, `
			INSERT INTO show_tiles (show_id, tile_id, weight, score)
			VALUES ($1, $2, $3, $4)
		`, showID, tileID, 1.0, 5.0)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, `
			INSERT INTO show_tiles (show_id, tile_id, weight, score)
			VALUES ($1, $2, $3, $4)
		`, showID, tileID, 1.0, 5.0)
		return err
	}
}

// PopulateShowTilesWithRandom selects tilesPerShow random tiles and associates them with the show
func PopulateShowTilesWithRandom(ctx context.Context, showID string, tx ...pgx.Tx) error {
	log.Printf("PopulateShowTilesWithRandom called for show %s", showID)
	var rows pgx.Rows
	var err error

	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, `
			SELECT id FROM tiles ORDER BY random() LIMIT $1
		`, tilesPerShow)
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		rows, err = pool.Query(ctx, `
			SELECT id FROM tiles ORDER BY random() LIMIT $1
		`, tilesPerShow)
	}

	if err != nil {
		log.Printf("Error querying tiles: %v", err)
		return err
	}
	defer rows.Close()

	var tileIDs []string
	for rows.Next() {
		var tileID string
		err := rows.Scan(&tileID)
		if err != nil {
			log.Printf("Error scanning tileID: %v", err)
			return err
		}
		tileIDs = append(tileIDs, tileID)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error with rows: %v", err)
		return err
	}

	log.Printf("Selected %d tiles for show %s", len(tileIDs), showID)

	// Insert into show_tiles
	for _, tileID := range tileIDs {
		err = EnsureTileInShowTiles(ctx, showID, tileID, tx...)
		if err != nil {
			log.Printf("Error inserting tile %s for show %s: %v", tileID, showID, err)
			return err
		}
	}

	log.Printf("Inserted %d tiles into show_tiles for show %s", len(tileIDs), showID)
	return nil
}
