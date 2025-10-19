package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"wanshow-bingo/avatar"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
	"github.com/matoous/go-nanoid/v2"
)

// PersistPlayer saves or updates a Player in the database
func PersistPlayer(ctx context.Context, player *models.Player, tx ...pgx.Tx) error {
	if len(tx) > 0 {
		// Use provided transaction
		if player.ID == "" {
			// New player, generate ID and insert
			player.ID, _ = gonanoid.New(10)
			_, err := tx[0].Exec(ctx, `
				INSERT INTO players (id, did, display_name, avatar, settings, score, permissions)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, player.ID, player.DID, player.DisplayName, player.Avatar, player.Settings, player.Score, uint64(player.Permissions))
			return err
		} else {
			// Existing player, update
			_, err := tx[0].Exec(ctx, `
				UPDATE players
				SET did = $1, display_name = $2, avatar = $3, settings = $4, score = $5, permissions = $6, updated_at = CURRENT_TIMESTAMP
				WHERE id = $7 AND deleted_at IS NULL
			`, player.DID, player.DisplayName, player.Avatar, player.Settings, player.Score, uint64(player.Permissions), player.ID)
			return err
		}
	} else {
		// Use pool directly
		pool := Pool()
		if pool == nil {
			return errors.New("database not available")
		}

		if player.ID == "" {
			// New player, generate ID and insert
			player.ID, _ = gonanoid.New(10)
			_, err := pool.Exec(ctx, `
				INSERT INTO players (id, did, display_name, avatar, settings, score, permissions)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, player.ID, player.DID, player.DisplayName, player.Avatar, player.Settings, player.Score, uint64(player.Permissions))
			return err
		} else {
			// Existing player, update
			_, err := pool.Exec(ctx, `
				UPDATE players
				SET did = $1, display_name = $2, avatar = $3, settings = $4, score = $5, permissions = $6, updated_at = CURRENT_TIMESTAMP
				WHERE id = $7 AND deleted_at IS NULL
			`, player.DID, player.DisplayName, player.Avatar, player.Settings, player.Score, uint64(player.Permissions), player.ID)
			return err
		}
	}
}

// GetPlayerByID retrieves a Player by ID
func GetPlayerByID(ctx context.Context, id string, tx ...pgx.Tx) (*models.Player, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at
			FROM players
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at
			FROM players
			WHERE id = $1 AND deleted_at IS NULL
		`, id)
	}

	var player models.Player
	err := row.Scan(
		&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score, &player.Permissions,
		&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("player not found")
		}
		return nil, err
	}

	return &player, nil
}

// PersistPlayerToTx saves or updates a Player using a transaction (deprecated, use PersistPlayer with tx param)
func PersistPlayerToTx(ctx context.Context, tx pgx.Tx, player *models.Player) error {
	return PersistPlayer(ctx, player, tx)
}

// LoadPlayerFromTx loads a Player by ID using a transaction (deprecated, use GetPlayerByID with tx param)
func LoadPlayerFromTx(ctx context.Context, tx pgx.Tx, id string) (*models.Player, error) {
	return GetPlayerByID(ctx, id, tx)
}

// FindOrCreatePlayer finds an existing player by Discord ID or creates a new one
func FindOrCreatePlayer(ctx context.Context, discordUser *models.DiscordUser, tx ...pgx.Tx) (*models.Player, error) {
	if len(tx) > 0 {
		// Use transaction
		// First try to find existing player
		var player models.Player
		err := tx[0].QueryRow(ctx, "SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at FROM players WHERE did = $1 AND deleted_at IS NULL", discordUser.ID).Scan(
			&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score, &player.Permissions,
			&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
		)

		if err == nil {
			// Player exists, update avatar and display_name if needed
			updates := []string{}
			args := []interface{}{}
			argCount := 1

			// Check if display_name needs update
			newDisplayName := discordUser.Username
			if discordUser.GlobalName != "" {
				newDisplayName = discordUser.GlobalName
			}
			if player.DisplayName != newDisplayName {
				updates = append(updates, "display_name = $"+fmt.Sprintf("%d", argCount))
				args = append(args, newDisplayName)
				argCount++
			}

			avatarURL, err := avatar.UploadDiscordAvatar(discordUser.ID, discordUser.Avatar, player.ID)
			if err != nil {
				log.Printf("database: failed to upload avatar: %v", err)
				// Continue without updating avatar
			} else if avatarURL != "" {
				updates = append(updates, "avatar = $"+fmt.Sprintf("%d", argCount))
				args = append(args, avatarURL)
				argCount++
			}

			if len(updates) > 0 {
				query := "UPDATE players SET " + strings.Join(updates, ", ") + ", updated_at = CURRENT_TIMESTAMP WHERE id = $" + fmt.Sprintf("%d", argCount)
				args = append(args, player.ID)
				_, err = tx[0].Exec(ctx, query, args...)
				if err != nil {
					log.Printf("database: failed to update player: %v", err)
				}
			}
			return &player, nil
		}

		if err != pgx.ErrNoRows {
			return nil, err
		}

		// Player doesn't exist, create new one
		playerID, _ := gonanoid.New(10)

		// TEMP: Use default permissions for all users
		defaultPerms := models.DefaultPermissions()
		log.Printf("PLAYER: Granting default permissions: %d", uint64(defaultPerms))

		// Upload avatar
		log.Printf("PLAYER: Attempting to upload avatar for new player %s, Discord user %s, avatar hash: %s", playerID, discordUser.ID, discordUser.Avatar)
		avatarURL, err := avatar.UploadDiscordAvatar(discordUser.ID, discordUser.Avatar, playerID)
		if err != nil {
			log.Printf("database: failed to upload avatar for new player: %v", err)
			avatarURL = "" // Continue without avatar
		} else {
			log.Printf("PLAYER: Successfully uploaded avatar, URL: %s", avatarURL)
		}

		displayName := discordUser.Username
		if discordUser.GlobalName != "" {
			displayName = discordUser.GlobalName
		}

		log.Printf("PLAYER: Inserting player with ID: %s, permissions: %d", playerID, uint64(defaultPerms))
		_, err = tx[0].Exec(ctx, `
			INSERT INTO players (id, did, display_name, avatar, settings, score, permissions)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, playerID, discordUser.ID, displayName, avatarURL, "{}", 0, uint64(defaultPerms))

		if err != nil {
			log.Printf("PLAYER: Failed to insert player: %v", err)
			return nil, err
		}
		log.Printf("PLAYER: Successfully inserted player %s", playerID)

		// Fetch the created player
		err = tx[0].QueryRow(ctx, "SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at FROM players WHERE id = $1", playerID).Scan(
			&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score, &player.Permissions,
			&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
		)

		return &player, err
	} else {
		// Use pool directly
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}

		// First try to find existing player
		var player models.Player
		err := pool.QueryRow(ctx, "SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at FROM players WHERE did = $1 AND deleted_at IS NULL", discordUser.ID).Scan(
			&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score, &player.Permissions,
			&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
		)

		if err == nil {
			// Player exists, update avatar and display_name if needed
			updates := []string{}
			args := []interface{}{}
			argCount := 1

			// Check if display_name needs update
			newDisplayName := discordUser.Username
			if discordUser.GlobalName != "" {
				newDisplayName = discordUser.GlobalName
			}
			if player.DisplayName != newDisplayName {
				updates = append(updates, "display_name = $"+fmt.Sprintf("%d", argCount))
				args = append(args, newDisplayName)
				argCount++
			}

			avatarURL, err := avatar.UploadDiscordAvatar(discordUser.ID, discordUser.Avatar, player.ID)
			if err != nil {
				log.Printf("database: failed to upload avatar: %v", err)
				// Continue without updating avatar
			} else if avatarURL != "" {
				updates = append(updates, "avatar = $"+fmt.Sprintf("%d", argCount))
				args = append(args, avatarURL)
				argCount++
			}

			if len(updates) > 0 {
				query := "UPDATE players SET " + strings.Join(updates, ", ") + ", updated_at = CURRENT_TIMESTAMP WHERE id = $" + fmt.Sprintf("%d", argCount)
				args = append(args, player.ID)
				_, err = pool.Exec(ctx, query, args...)
				if err != nil {
					log.Printf("database: failed to update player: %v", err)
				}
			}
			return &player, nil
		}

		if err != pgx.ErrNoRows {
			return nil, err
		}

		// Player doesn't exist, create new one
		playerID, _ := gonanoid.New(10)

		// Check if this is an admin user
		adminDiscordIDs := []string{
			"180639594017062912", // altrius_codes
			"639238547999817769",
			"773038992727343164",
			"171160105155297282",
			"186845156404756480",
		}

		var defaultPerms models.Permission
		isAdmin := false
		for _, adminID := range adminDiscordIDs {
			if discordUser.ID == adminID {
				defaultPerms = models.AdminPermissions()
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			defaultPerms = models.DefaultPermissions()
		}

		// Upload avatar
		avatarURL, err := avatar.UploadDiscordAvatar(discordUser.ID, discordUser.Avatar, playerID)
		if err != nil {
			log.Printf("database: failed to upload avatar for new player: %v", err)
			avatarURL = "" // Continue without avatar
		}

		displayName := discordUser.Username
		if discordUser.GlobalName != "" {
			displayName = discordUser.GlobalName
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO players (id, did, display_name, avatar, settings, score, permissions)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, playerID, discordUser.ID, displayName, avatarURL, "{}", 0, uint64(defaultPerms))

		if err != nil {
			return nil, err
		}

		// Fetch the created player
		err = pool.QueryRow(ctx, "SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at FROM players WHERE id = $1", playerID).Scan(
			&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score, &player.Permissions,
			&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
		)

		return &player, err
	}
}

// GetPlayerByIdentifier finds a player by ID or display_name
func GetPlayerByIdentifier(ctx context.Context, identifier string, tx ...pgx.Tx) (*models.Player, error) {
	var row pgx.Row

	if len(tx) > 0 {
		row = tx[0].QueryRow(ctx, `
			SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at
			FROM players
			WHERE (id = $1 OR display_name = $1) AND deleted_at IS NULL
		`, identifier)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		row = pool.QueryRow(ctx, `
			SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at
			FROM players
			WHERE (id = $1 OR display_name = $1) AND deleted_at IS NULL
		`, identifier)
	}

	var player models.Player
	err := row.Scan(
		&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score, &player.Permissions,
		&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("player not found")
		}
		return nil, err
	}

	return &player, nil
}

// GetAllPlayers returns all players (for admin purposes, might want to add pagination later)
func GetAllPlayers(ctx context.Context, tx ...pgx.Tx) ([]models.Player, error) {
	var rows pgx.Rows
	var err error

	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, `
			SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at
			FROM players
			WHERE deleted_at IS NULL
			ORDER BY display_name
		`)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, errors.New("database not available")
		}
		rows, err = pool.Query(ctx, `
			SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at
			FROM players
			WHERE deleted_at IS NULL
			ORDER BY display_name
		`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		err := rows.Scan(
			&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score, &player.Permissions,
			&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, rows.Err()
}

// GetAllPlayersPaginated returns all players with pagination support
func GetAllPlayersPaginated(ctx context.Context, limit, offset int, orderBy, orderDir string, tx ...pgx.Tx) ([]models.Player, int, error) {
	var rows pgx.Rows
	var countRow pgx.Row
	var err error

	// Validate order parameters for security
	validOrderBy := map[string]bool{
		"display_name": true,
		"created_at":   true,
		"score":        true,
	}
	if !validOrderBy[orderBy] {
		orderBy = "display_name"
	}
	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "asc"
	}

	query := `
		SELECT id, did, display_name, avatar, settings, score, permissions, created_at, updated_at, deleted_at
		FROM players
		WHERE deleted_at IS NULL
		ORDER BY ` + orderBy + ` ` + orderDir + `
		LIMIT $1 OFFSET $2
	`

	countQuery := `SELECT COUNT(*) FROM players WHERE deleted_at IS NULL`

	if len(tx) > 0 {
		rows, err = tx[0].Query(ctx, query, limit, offset)
		countRow = tx[0].QueryRow(ctx, countQuery)
	} else {
		pool := Pool()
		if pool == nil {
			return nil, 0, errors.New("database not available")
		}
		rows, err = pool.Query(ctx, query, limit, offset)
		countRow = pool.QueryRow(ctx, countQuery)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		err := rows.Scan(
			&player.ID, &player.DID, &player.DisplayName, &player.Avatar, &player.Settings, &player.Score, &player.Permissions,
			&player.CreatedAt, &player.UpdatedAt, &player.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var totalCount int
	err = countRow.Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return players, totalCount, nil
}
