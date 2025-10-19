package db

import (
	"context"
	"errors"
	"time"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
	"github.com/matoous/go-nanoid/v2"
)

// PersistTimer saves or updates a Timer in the database
func PersistTimer(ctx context.Context, timer *models.Timer, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		// Use transaction
		if timer.ID == "" {
			// New timer, generate ID and insert
			timer.ID, _ = gonanoid.New(10)
			_, err := tx[0].Exec(ctx, `
				INSERT INTO timers (id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, timer.ID, timer.Title, timer.Duration, timer.CreatedBy, timer.ShowID, timer.StartsAt, timer.ExpiresAt, timer.IsActive, timer.Settings)
			return err
		} else {
			// Existing timer, update
			_, err := tx[0].Exec(ctx, `
				UPDATE timers
				SET title = $1, duration = $2, created_by = $3, show_id = $4, starts_at = $5, expires_at = $6, is_active = $7, settings = $8, updated_at = CURRENT_TIMESTAMP
				WHERE id = $9 AND deleted_at IS NULL
			`, timer.Title, timer.Duration, timer.CreatedBy, timer.ShowID, timer.StartsAt, timer.ExpiresAt, timer.IsActive, timer.Settings, timer.ID)
			return err
		}
	} else {
		// Use pool directly
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}

		if timer.ID == "" {
			// New timer, generate ID and insert
			timer.ID, _ = gonanoid.New(10)
			_, err := pool.Exec(ctx, `
				INSERT INTO timers (id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, timer.ID, timer.Title, timer.Duration, timer.CreatedBy, timer.ShowID, timer.StartsAt, timer.ExpiresAt, timer.IsActive, timer.Settings)
			return err
		} else {
			// Existing timer, update
			_, err := pool.Exec(ctx, `
				UPDATE timers
				SET title = $1, duration = $2, created_by = $3, show_id = $4, starts_at = $5, expires_at = $6, is_active = $7, settings = $8, updated_at = CURRENT_TIMESTAMP
				WHERE id = $9 AND deleted_at IS NULL
			`, timer.Title, timer.Duration, timer.CreatedBy, timer.ShowID, timer.StartsAt, timer.ExpiresAt, timer.IsActive, timer.Settings, timer.ID)
			return err
		}
	}
}

// GetTimerByID retrieves a Timer by ID
func GetTimerByID(ctx context.Context, id string, tx ...pgx.Tx) (*models.Timer, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings, created_at, updated_at, deleted_at
			FROM timers
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings, created_at, updated_at, deleted_at
			FROM timers
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	}

	var timer models.Timer
	err := row.Scan(
		&timer.ID, &timer.Title, &timer.Duration, &timer.CreatedBy, &timer.ShowID,
		&timer.StartsAt, &timer.ExpiresAt, &timer.IsActive, &timer.Settings,
		&timer.CreatedAt, &timer.UpdatedAt, &timer.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("timer not found")
		}
		return nil, err
	}

	return &timer, nil
}

// GetActiveTimersForShow retrieves all active timers for a specific show
func GetActiveTimersForShow(ctx context.Context, showID string, tx ...pgx.Tx) ([]models.Timer, error) {
	var rows pgx.Rows
	var err error

	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, `
			SELECT id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings, created_at, updated_at, deleted_at
			FROM timers
			WHERE show_id = $1 AND is_active = true AND deleted_at IS NULL
			ORDER BY created_at DESC
		`, showID)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		rows, err = pool.Query(ctx, `
			SELECT id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings, created_at, updated_at, deleted_at
			FROM timers
			WHERE show_id = $1 AND is_active = true AND deleted_at IS NULL
			ORDER BY created_at DESC
		`, showID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timers []models.Timer
	for rows.Next() {
		var timer models.Timer
		err := rows.Scan(
			&timer.ID, &timer.Title, &timer.Duration, &timer.CreatedBy, &timer.ShowID,
			&timer.StartsAt, &timer.ExpiresAt, &timer.IsActive, &timer.Settings,
			&timer.CreatedAt, &timer.UpdatedAt, &timer.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		timers = append(timers, timer)
	}

	return timers, rows.Err()
}

// StartTimer activates a timer and sets its expiration time
func StartTimer(ctx context.Context, timerID string, tx ...pgx.Tx) error {
	now := time.Now()

	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, `
			UPDATE timers
			SET is_active = true, starts_at = $1, expires_at = $1 + make_interval(secs => duration), updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND deleted_at IS NULL
		`, now, timerID)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, `
			UPDATE timers
			SET is_active = true, starts_at = $1, expires_at = $1 + make_interval(secs => duration), updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND deleted_at IS NULL
		`, now, timerID)
		return err
	}
}

// StopTimer deactivates a timer
func StopTimer(ctx context.Context, timerID string, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, `
			UPDATE timers
			SET is_active = false, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND deleted_at IS NULL
		`, timerID)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, `
			UPDATE timers
			SET is_active = false, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND deleted_at IS NULL
		`, timerID)
		return err
	}
}

// StopActiveTimersByTitle deactivates all active timers with the given title for a show
func StopActiveTimersByTitle(ctx context.Context, title string, showID string) error {
	pool := Pool()
	if pool == nil {
		return errors.New("database not available")
	}
	_, err := pool.Exec(ctx, `
		UPDATE timers
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE title = $1 AND show_id = $2 AND is_active = true AND deleted_at IS NULL
	`, title, showID)
	return err
}

// GetExpiredTimers retrieves timers that have expired
func GetExpiredTimers(ctx context.Context, tx ...pgx.Tx) ([]models.Timer, error) {
	var rows pgx.Rows
	var err error

	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, `
			SELECT id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings, created_at, updated_at, deleted_at
			FROM timers
			WHERE is_active = true AND expires_at < CURRENT_TIMESTAMP AND deleted_at IS NULL
		`)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		rows, err = pool.Query(ctx, `
			SELECT id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings, created_at, updated_at, deleted_at
			FROM timers
			WHERE is_active = true AND expires_at < CURRENT_TIMESTAMP AND deleted_at IS NULL
		`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timers []models.Timer
	for rows.Next() {
		var timer models.Timer
		err := rows.Scan(
			&timer.ID, &timer.Title, &timer.Duration, &timer.CreatedBy, &timer.ShowID,
			&timer.StartsAt, &timer.ExpiresAt, &timer.IsActive, &timer.Settings,
			&timer.CreatedAt, &timer.UpdatedAt, &timer.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		timers = append(timers, timer)
	}

	return timers, rows.Err()
}

// DeleteTimer soft deletes a timer
func DeleteTimer(ctx context.Context, timerID string, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, "UPDATE timers SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", timerID)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, "UPDATE timers SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", timerID)
		return err
	}
}
