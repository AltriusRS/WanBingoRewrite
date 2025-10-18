package db

import (
	"context"
	"time"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
	"github.com/matoous/go-nanoid/v2"
)

// CreateTileSuggestion inserts a new tile suggestion into the database
func CreateTileSuggestion(ctx context.Context, name, tileName, reason string) (*models.TileSuggestion, error) {
	id, _ := gonanoid.New(10)
	now := time.Now()

	query := `
		INSERT INTO tile_suggestions (id, name, tile_name, reason, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'pending', $5, $6)
		RETURNING id, name, tile_name, reason, status, reviewed_by, reviewed_at, created_at, updated_at, deleted_at
	`

	var suggestion models.TileSuggestion
	err := Pool().QueryRow(ctx, query, id, name, tileName, reason, now, now).Scan(
		&suggestion.ID,
		&suggestion.Name,
		&suggestion.TileName,
		&suggestion.Reason,
		&suggestion.Status,
		&suggestion.ReviewedBy,
		&suggestion.ReviewedAt,
		&suggestion.CreatedAt,
		&suggestion.UpdatedAt,
		&suggestion.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &suggestion, nil
}

// GetTileSuggestions retrieves all tile suggestions, optionally filtered by status
func GetTileSuggestions(ctx context.Context, status *string) ([]models.TileSuggestion, error) {
	var query string
	var args []interface{}

	if status != nil {
		query = `SELECT id, name, tile_name, reason, status, reviewed_by, reviewed_at, created_at, updated_at, deleted_at
				 FROM tile_suggestions
				 WHERE status = $1 AND deleted_at IS NULL
				 ORDER BY created_at DESC`
		args = []interface{}{*status}
	} else {
		query = `SELECT id, name, tile_name, reason, status, reviewed_by, reviewed_at, created_at, updated_at, deleted_at
				 FROM tile_suggestions
				 WHERE deleted_at IS NULL
				 ORDER BY created_at DESC`
	}

	rows, err := Pool().Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []models.TileSuggestion
	for rows.Next() {
		var suggestion models.TileSuggestion
		err := rows.Scan(
			&suggestion.ID,
			&suggestion.Name,
			&suggestion.TileName,
			&suggestion.Reason,
			&suggestion.Status,
			&suggestion.ReviewedBy,
			&suggestion.ReviewedAt,
			&suggestion.CreatedAt,
			&suggestion.UpdatedAt,
			&suggestion.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// UpdateTileSuggestion updates a tile suggestion's status and review info
func UpdateTileSuggestion(ctx context.Context, id, status string, reviewedBy *string) (*models.TileSuggestion, error) {
	now := time.Now()
	var reviewedAt *time.Time
	if reviewedBy != nil {
		reviewedAt = &now
	}

	query := `
		UPDATE tile_suggestions
		SET status = $1, reviewed_by = $2, reviewed_at = $3, updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id, name, tile_name, reason, status, reviewed_by, reviewed_at, created_at, updated_at, deleted_at
	`

	var suggestion models.TileSuggestion
	err := Pool().QueryRow(ctx, query, status, reviewedBy, reviewedAt, now, id).Scan(
		&suggestion.ID,
		&suggestion.Name,
		&suggestion.TileName,
		&suggestion.Reason,
		&suggestion.Status,
		&suggestion.ReviewedBy,
		&suggestion.ReviewedAt,
		&suggestion.CreatedAt,
		&suggestion.UpdatedAt,
		&suggestion.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &suggestion, nil
}
