package db

import (
	"context"
	"errors"
	"time"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
)

// PersistSession saves or updates a Session in the database
func PersistSession(ctx context.Context, session *models.Session, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		// Use transaction
		if session.ID == "" {
			// New session, generate ID and insert
			session.ID = generateSessionID()
			if session.ExpiresAt.IsZero() {
				session.ExpiresAt = time.Now().Add(24 * time.Hour)
			}
			_, err := tx[0].Exec(ctx, `
				INSERT INTO sessions (id, player_id, expires_at)
				VALUES ($1, $2, $3)
			`, session.ID, session.PlayerID, session.ExpiresAt)
			return err
		} else {
			// Existing session, update
			_, err := tx[0].Exec(ctx, `
				UPDATE sessions
				SET player_id = $1, expires_at = $2
				WHERE id = $3 AND deleted_at IS NULL
			`, session.PlayerID, session.ExpiresAt, session.ID)
			return err
		}
	} else {
		// Use pool directly
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}

		if session.ID == "" {
			// New session, generate ID and insert
			session.ID = generateSessionID()
			if session.ExpiresAt.IsZero() {
				session.ExpiresAt = time.Now().Add(24 * time.Hour)
			}
			_, err := pool.Exec(ctx, `
				INSERT INTO sessions (id, player_id, expires_at)
				VALUES ($1, $2, $3)
			`, session.ID, session.PlayerID, session.ExpiresAt)
			return err
		} else {
			// Existing session, update
			_, err := pool.Exec(ctx, `
				UPDATE sessions
				SET player_id = $1, expires_at = $2
				WHERE id = $3 AND deleted_at IS NULL
			`, session.PlayerID, session.ExpiresAt, session.ID)
			return err
		}
	}
}

// GetSessionByID retrieves a Session by ID
func GetSessionByID(ctx context.Context, id string, tx ...pgx.Tx) (*models.Session, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, player_id, created_at, expires_at, deleted_at
			FROM sessions
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, player_id, created_at, expires_at, deleted_at
			FROM sessions
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	}

	var session models.Session
	err := row.Scan(
		&session.ID, &session.PlayerID, &session.CreatedAt, &session.ExpiresAt, &session.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	return &session, nil
}

// PersistSessionToTx saves or updates a Session using a transaction (deprecated, use PersistSession with tx param)
func PersistSessionToTx(ctx context.Context, tx pgx.Tx, session *models.Session) error {
	return PersistSession(ctx, session, tx)
}

// LoadSessionFromTx loads a Session by ID using a transaction (deprecated, use GetSessionByID with tx param)
func LoadSessionFromTx(ctx context.Context, tx pgx.Tx, id string) (*models.Session, error) {
	return GetSessionByID(ctx, id, tx)
}

// CreateSession creates a new session for a player
func CreateSession(ctx context.Context, playerID string, tx ...pgx.Tx) (string, error) {
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(24 * time.Hour) // Sessions last 24 hours

	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, `
			INSERT INTO sessions (id, player_id, expires_at)
			VALUES ($1, $2, $3)
		`, sessionID, playerID, expiresAt)
		return sessionID, err
	} else {
		pool := Pool()
		if pool == nil {
			return "", errors.New("database not available")
		}
		_, err := pool.Exec(ctx, `
			INSERT INTO sessions (id, player_id, expires_at)
			VALUES ($1, $2, $3)
		`, sessionID, playerID, expiresAt)
		return sessionID, err
	}
}

// ValidateSession checks if a session is valid and returns the player
func ValidateSession(ctx context.Context, sessionID string, tx ...pgx.Tx) (*models.Player, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT p.* FROM players p
			JOIN sessions s ON p.id = s.player_id
			WHERE s.id = $1 AND s.expires_at > CURRENT_TIMESTAMP AND s.deleted_at IS NULL AND p.deleted_at IS NULL
		`, sessionID)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT p.* FROM players p
			JOIN sessions s ON p.id = s.player_id
			WHERE s.id = $1 AND s.expires_at > CURRENT_TIMESTAMP AND s.deleted_at IS NULL AND p.deleted_at IS NULL
		`, sessionID)
	}

	var player models.Player
	err := row.Scan(
		&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score,
		&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("invalid session")
		}
		return nil, err
	}

	return &player, nil
}

// DeleteSession invalidates a session
func DeleteSession(ctx context.Context, sessionID string, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, "UPDATE sessions SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", sessionID)
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, "UPDATE sessions SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", sessionID)
		return err
	}
}

// CleanupExpiredSessions removes expired sessions from the database
func CleanupExpiredSessions(ctx context.Context, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		_, err := tx[0].Exec(ctx, "DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP")
		return err
	} else {
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}
		_, err := pool.Exec(ctx, "DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP")
		return err
	}
}
