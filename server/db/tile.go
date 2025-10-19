package db

import (
	"context"
	"errors"
	"log"
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

// PersistTile saves or updates a Tile in the database
func PersistTile(ctx context.Context, tile *models.Tile, tx ...pgx.Tx) error {
	log.Printf("PersistTile called for tile %s with settings: %+v", tile.ID, tile.Settings)
	if len(tx) > 0 {
		result, err := tx[0].Exec(ctx, `
			INSERT INTO tiles (id, title, category, weight, score, created_by, settings, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id) DO UPDATE SET
				title = EXCLUDED.title,
				category = EXCLUDED.category,
				weight = EXCLUDED.weight,
				score = EXCLUDED.score,
				settings = EXCLUDED.settings,
				updated_at = EXCLUDED.updated_at
		`, tile.ID, tile.Title, tile.Category, tile.Weight, tile.Score, tile.CreatedBy, tile.Settings, tile.CreatedAt, tile.UpdatedAt)
		log.Printf("PersistTile result: %v, err: %v", result, err)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		result, err := pool.Exec(ctx, `
			INSERT INTO tiles (id, title, category, weight, score, created_by, settings, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id) DO UPDATE SET
				title = EXCLUDED.title,
				category = EXCLUDED.category,
				weight = EXCLUDED.weight,
				score = EXCLUDED.score,
				settings = EXCLUDED.settings,
				updated_at = EXCLUDED.updated_at
		`, tile.ID, tile.Title, tile.Category, tile.Weight, tile.Score, tile.CreatedBy, tile.Settings, tile.CreatedAt, tile.UpdatedAt)
		log.Printf("PersistTile result: %v, err: %v", result, err)
		return err
	}
}

// DeleteTile soft deletes a tile by setting deleted_at
func DeleteTile(ctx context.Context, id string, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, `
			UPDATE tiles SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
		`, id)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, `
			UPDATE tiles SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
		`, id)
		return err
	}
}

// DeleteTileConfirmation removes a tile confirmation for a specific show and tile
func DeleteTileConfirmation(ctx context.Context, showID, tileID string, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, `
			UPDATE tile_confirmations SET deleted_at = CURRENT_TIMESTAMP
			WHERE show_id = $1 AND tile_id = $2 AND deleted_at IS NULL
		`, showID, tileID)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, `
			UPDATE tile_confirmations SET deleted_at = CURRENT_TIMESTAMP
			WHERE show_id = $1 AND tile_id = $2 AND deleted_at IS NULL
		`, showID, tileID)
		return err
	}
}

// GetTileConfirmationsForShow retrieves all tile confirmations for a specific show
// Only returns confirmations for tiles that still exist (not deleted)
func GetTileConfirmationsForShow(ctx context.Context, showID string, tx ...pgx.Tx) ([]models.TileConfirmation, error) {
	var rows pgx.Rows
	var err error

	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, `
			SELECT tc.id, tc.show_id, tc.tile_id, tc.confirmed_by, tc.context, tc.confirmation_time, tc.created_at, tc.updated_at, tc.deleted_at
			FROM tile_confirmations tc
			INNER JOIN tiles t ON tc.tile_id = t.id
			WHERE tc.show_id = $1 AND tc.deleted_at IS NULL AND t.deleted_at IS NULL
			ORDER BY tc.confirmation_time
		`, showID)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		rows, err = pool.Query(ctx, `
			SELECT tc.id, tc.show_id, tc.tile_id, tc.confirmed_by, tc.context, tc.confirmation_time, tc.created_at, tc.updated_at, tc.deleted_at
			FROM tile_confirmations tc
			INNER JOIN tiles t ON tc.tile_id = t.id
			WHERE tc.show_id = $1 AND tc.deleted_at IS NULL AND t.deleted_at IS NULL
			ORDER BY tc.confirmation_time
		`, showID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var confirmations []models.TileConfirmation
	for rows.Next() {
		var conf models.TileConfirmation
		err := rows.Scan(
			&conf.ID,
			&conf.ShowID,
			&conf.TileID,
			&conf.ConfirmedBy,
			&conf.Context,
			&conf.ConfirmationTime,
			&conf.CreatedAt,
			&conf.UpdatedAt,
			&conf.DeletedAt,
		)
		if err != nil {
			log.Printf("Error scanning tile confirmation row: %v", err)
			return nil, err
		}
		confirmations = append(confirmations, conf)
	}

	return confirmations, nil
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
