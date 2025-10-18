package tilerouter

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db/models"

	"wanshow-bingo/db"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

// Get retrieves tiles with optional pagination and filtering
func Get(c *fiber.Ctx) error {
	pool := db.Pool()

	// If no database is available, return an error
	if pool == nil {
		return utils.NewApiError("Failed to connect to database", 0x0101).AsResponse(c)
	}

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
	category := c.Query("category")
	orderBy := c.Query("order_by", "created_at")
	orderDir := c.Query("order_dir", "desc")

	// Validate order_by parameter for security
	validOrderBy := map[string]bool{
		"created_at": true,
		"title":      true,
		"category":   true,
		"weight":     true,
		"score":      true,
	}
	if !validOrderBy[orderBy] {
		orderBy = "created_at"
	}

	// Validate order direction
	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "desc"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Build query with optional category filter
	var query string
	var args []interface{}
	var countQuery string
	var countArgs []interface{}

	if category != "" {
		query = `
			SELECT id, title, category, last_drawn, weight, score, created_by, settings, created_at, updated_at, deleted_at
			FROM tiles
			WHERE deleted_at IS NULL AND category = $1
			ORDER BY ` + orderBy + ` ` + orderDir + `
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{category, limit, offset}

		countQuery = `SELECT COUNT(*) FROM tiles WHERE deleted_at IS NULL AND category = $1`
		countArgs = []interface{}{category}
	} else {
		query = `
			SELECT id, title, category, last_drawn, weight, score, created_by, settings, created_at, updated_at, deleted_at
			FROM tiles
			WHERE deleted_at IS NULL
			ORDER BY ` + orderBy + ` ` + orderDir + `
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}

		countQuery = `SELECT COUNT(*) FROM tiles WHERE deleted_at IS NULL`
		countArgs = []interface{}{}
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("tiles query error: %v", err)
		return utils.NewApiError("Failed to query tiles", 0x0201).AsResponse(c)
	}
	defer rows.Close()

	tiles, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Tile])
	if err != nil {
		log.Printf("tiles collection error: %v", err)
		return utils.NewApiError("Failed to process tiles", 0x0202).AsResponse(c)
	}

	// Get total count for pagination metadata
	var totalCount int
	err = pool.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		log.Printf("count query error: %v", err)
		// Don't fail the request for count error, just omit pagination info
		totalCount = 0
	}

	totalPages := (totalCount + limit - 1) / limit // Ceiling division

	return c.JSON(fiber.Map{
		"tiles": tiles,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_count": totalCount,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
		"filters": fiber.Map{
			"category":  category,
			"order_by":  orderBy,
			"order_dir": orderDir,
		},
	})
}
