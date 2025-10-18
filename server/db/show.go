package db

import (
	"context"
	"errors"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
)

// PersistShow saves or updates a Show in the database
func PersistShow(ctx context.Context, show *models.Show, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		// Use transaction
		if show.ID == "" {
			// New show, generate ID and insert
			show.ID = generateID(10)
			_, err := tx[0].Exec(ctx, `
				INSERT INTO shows (id, youtube_id, scheduled_time, actual_start_time, thumbnail, metadata)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, show.ID, show.YoutubeID, show.ScheduledTime, show.ActualStartTime, show.Thumbnail, show.Metadata)
			return err
		} else {
			// Existing show, update
			_, err := tx[0].Exec(ctx, `
				UPDATE shows
				SET youtube_id = $1, scheduled_time = $2, actual_start_time = $3, thumbnail = $4, metadata = $5, updated_at = CURRENT_TIMESTAMP
				WHERE id = $6 AND deleted_at IS NULL
			`, show.YoutubeID, show.ScheduledTime, show.ActualStartTime, show.Thumbnail, show.Metadata, show.ID)
			return err
		}
	} else {
		// Use pool directly
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}

		if show.ID == "" {
			// New show, generate ID and insert
			show.ID = generateID(10)
			_, err := pool.Exec(ctx, `
				INSERT INTO shows (id, youtube_id, scheduled_time, actual_start_time, thumbnail, metadata)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, show.ID, show.YoutubeID, show.ScheduledTime, show.ActualStartTime, show.Thumbnail, show.Metadata)
			return err
		} else {
			// Existing show, update
			_, err := pool.Exec(ctx, `
				UPDATE shows
				SET youtube_id = $1, scheduled_time = $2, actual_start_time = $3, thumbnail = $4, metadata = $5, updated_at = CURRENT_TIMESTAMP
				WHERE id = $6 AND deleted_at IS NULL
			`, show.YoutubeID, show.ScheduledTime, show.ActualStartTime, show.Thumbnail, show.Metadata, show.ID)
			return err
		}
	}
}

// GetShowByID retrieves a Show by ID
func GetShowByID(ctx context.Context, id string, tx ...pgx.Tx) (*models.Show, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, youtube_id, scheduled_time, actual_start_time, thumbnail, metadata, created_at, updated_at, deleted_at
			FROM shows
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, youtube_id, scheduled_time, actual_start_time, thumbnail, metadata, created_at, updated_at, deleted_at
			FROM shows
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	}

	var show models.Show
	err := row.Scan(
		&show.ID, &show.YoutubeID, &show.ScheduledTime, &show.ActualStartTime, &show.Thumbnail, &show.Metadata,
		&show.CreatedAt, &show.UpdatedAt, &show.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("show not found")
		}
		return nil, err
	}

	return &show, nil
}

// PersistShowToTx saves or updates a Show using a transaction (deprecated, use PersistShow with tx param)
func PersistShowToTx(ctx context.Context, tx pgx.Tx, show *models.Show) error {
	return PersistShow(ctx, show, tx)
}

// LoadShowFromTx loads a Show by ID using a transaction (deprecated, use GetShowByID with tx param)
func LoadShowFromTx(ctx context.Context, tx pgx.Tx, id string) (*models.Show, error) {
	return GetShowByID(ctx, id, tx)
}

// GetLatestShow retrieves the most recent show by scheduled time
func GetLatestShow(ctx context.Context, tx ...pgx.Tx) (*models.Show, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, youtube_id, scheduled_time, actual_start_time, thumbnail, metadata, created_at, updated_at, deleted_at
			FROM shows
			WHERE deleted_at IS NULL
			ORDER BY scheduled_time DESC
			LIMIT 1
		`)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, youtube_id, scheduled_time, actual_start_time, thumbnail, metadata, created_at, updated_at, deleted_at
			FROM shows
			WHERE deleted_at IS NULL
			ORDER BY scheduled_time DESC
			LIMIT 1
		`)
	}

	var show models.Show
	err := row.Scan(
		&show.ID, &show.YoutubeID, &show.ScheduledTime, &show.ActualStartTime, &show.Thumbnail, &show.Metadata,
		&show.CreatedAt, &show.UpdatedAt, &show.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("no shows found")
		}
		return nil, err
	}

	return &show, nil
}
