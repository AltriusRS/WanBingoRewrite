package users

import (
	"context"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// GetByIdentifier returns partial player profile by ID or display_name
func GetByIdentifier(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	if identifier == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Identifier parameter required", 400))
	}

	player, err := db.GetPlayerByIdentifier(context.Background(), identifier)
	if err != nil {
		if err.Error() == "player not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.NewApiError("Player not found", 404))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to fetch player", 500))
	}

	// Return partial profile (exclude sensitive info like DID)
	return c.JSON(fiber.Map{
		"success": true,
		"player": fiber.Map{
			"id":           player.ID,
			"display_name": player.DisplayName,
			"avatar":       player.Avatar,
			"score":        player.Score,
			"created_at":   player.CreatedAt,
		},
	})
}

// GetAll returns all players with pagination (partial profiles)
func GetAll(c *fiber.Ctx) error {
	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 50)
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	// Parse filter parameters
	orderBy := c.Query("order_by", "display_name")
	orderDir := c.Query("order_dir", "asc")

	// Validate order_by parameter for security
	validOrderBy := map[string]bool{
		"display_name": true,
		"created_at":   true,
		"score":        true,
	}
	if !validOrderBy[orderBy] {
		orderBy = "display_name"
	}

	// Validate order direction
	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "asc"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get paginated players from database
	players, totalCount, err := db.GetAllPlayersPaginated(ctx, limit, offset, orderBy, orderDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to fetch players", 500))
	}

	// Return partial profiles for players
	var playerProfiles []fiber.Map
	for _, player := range players {
		playerProfiles = append(playerProfiles, fiber.Map{
			"id":           player.ID,
			"display_name": player.DisplayName,
			"avatar":       player.Avatar,
			"score":        player.Score,
			"created_at":   player.CreatedAt,
		})
	}

	totalPages := (totalCount + limit - 1) / limit // Ceiling division

	return c.JSON(fiber.Map{
		"success": true,
		"players": playerProfiles,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_count": totalCount,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
		"filters": fiber.Map{
			"order_by":  orderBy,
			"order_dir": orderDir,
		},
	})
}
