package db

import (
	"context"
	"errors"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
)

// PersistMessage saves or updates a Message in the database
func PersistMessage(ctx context.Context, message *models.Message, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		// Use transaction
		if message.ID == "" {
			// New message, generate ID and insert
			message.ID = generateID(10)
			_, err := tx[0].Exec(ctx, `
				INSERT INTO messages (id, show_id, player_id, contents, system, replying)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, message.ID, message.ShowID, message.PlayerID, message.Contents, message.System, message.Replying)
			return err
		} else {
			// Existing message, update
			_, err := tx[0].Exec(ctx, `
				UPDATE messages
				SET show_id = $1, player_id = $2, contents = $3, system = $4, replying = $5, updated_at = CURRENT_TIMESTAMP
				WHERE id = $6 AND deleted_at IS NULL
			`, message.ShowID, message.PlayerID, message.Contents, message.System, message.Replying, message.ID)
			return err
		}
	} else {
		// Use pool directly
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}

		if message.ID == "" {
			// New message, generate ID and insert
			message.ID = generateID(10)
			_, err := pool.Exec(ctx, `
				INSERT INTO messages (id, show_id, player_id, contents, system, replying)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, message.ID, message.ShowID, message.PlayerID, message.Contents, message.System, message.Replying)
			return err
		} else {
			// Existing message, update
			_, err := pool.Exec(ctx, `
				UPDATE messages
				SET show_id = $1, player_id = $2, contents = $3, system = $4, replying = $5, updated_at = CURRENT_TIMESTAMP
				WHERE id = $6 AND deleted_at IS NULL
			`, message.ShowID, message.PlayerID, message.Contents, message.System, message.Replying, message.ID)
			return err
		}
	}
}

// GetMessageByID retrieves a Message by ID
func GetMessageByID(ctx context.Context, id string, tx ...pgx.Tx) (*models.Message, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, show_id, player_id, contents, system, replying, created_at, updated_at, deleted_at
			FROM messages
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, show_id, player_id, contents, system, replying, created_at, updated_at, deleted_at
			FROM messages
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	}

	var message models.Message
	err := row.Scan(
		&message.ID, &message.ShowID, &message.PlayerID, &message.Contents, &message.System, &message.Replying,
		&message.CreatedAt, &message.UpdatedAt, &message.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("message not found")
		}
		return nil, err
	}

	return &message, nil
}

// PersistMessageToTx saves or updates a Message using a transaction (deprecated, use PersistMessage with tx param)
func PersistMessageToTx(ctx context.Context, tx pgx.Tx, message *models.Message) error {
	return PersistMessage(ctx, message, tx)
}

// LoadMessageFromTx loads a Message by ID using a transaction (deprecated, use GetMessageByID with tx param)
func LoadMessageFromTx(ctx context.Context, tx pgx.Tx, id string) (*models.Message, error) {
	return GetMessageByID(ctx, id, tx)
}

// SaveMessage saves a message for the latest show (legacy function for compatibility)
func SaveMessage(ctx context.Context, msg *models.Message, tx ...pgx.Tx) error {
	latestShow, err := GetLatestShow(ctx, tx...)
	if err != nil {
		return err
	}

	msg.ShowID = latestShow.ID
	return PersistMessage(ctx, msg, tx...)
}

// GetMessageHistory retrieves recent messages for the latest show
func GetMessageHistory(ctx context.Context, tx ...pgx.Tx) ([]models.Message, error) {
	latestShow, err := GetLatestShow(ctx, tx...)
	if err != nil {
		return nil, err
	}

	var rows pgx.Rows
	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, `
			SELECT id, show_id, player_id, contents, system, replying, created_at, updated_at, deleted_at
			FROM messages
			WHERE show_id = $1 AND deleted_at IS NULL
			ORDER BY created_at DESC
			LIMIT 30
		`, latestShow.ID)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		rows, err = pool.Query(ctx, `
			SELECT id, show_id, player_id, contents, system, replying, created_at, updated_at, deleted_at
			FROM messages
			WHERE show_id = $1 AND deleted_at IS NULL
			ORDER BY created_at DESC
			LIMIT 30
		`, latestShow.ID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID, &message.ShowID, &message.PlayerID, &message.Contents, &message.System, &message.Replying,
			&message.CreatedAt, &message.UpdatedAt, &message.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, rows.Err()
}
