package timers

import (
	"context"
	"log"
	"strings"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

// GetTimers retrieves timers with optional pagination and filtering
func GetTimers(c *fiber.Ctx) error {
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
	showID := c.Query("show_id")
	isActive := c.Query("is_active")
	orderBy := c.Query("order_by", "created_at")
	orderDir := c.Query("order_dir", "desc")

	// Validate order_by parameter for security
	validOrderBy := map[string]bool{
		"created_at": true,
		"title":      true,
		"duration":   true,
		"expires_at": true,
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

	// Build query with optional filters
	var query string
	var args []interface{}
	var countQuery string
	var countArgs []interface{}

	baseQuery := `
		SELECT id, title, duration, created_by, show_id, starts_at, expires_at, is_active, settings, created_at, updated_at, deleted_at
		FROM timers
		WHERE deleted_at IS NULL
	`
	baseCountQuery := `SELECT COUNT(*) FROM timers WHERE deleted_at IS NULL`

	conditions := []string{}
	countConditions := []string{}

	if showID != "" {
		conditions = append(conditions, "show_id = $"+string(rune(len(args)+1)))
		countConditions = append(countConditions, "show_id = $"+string(rune(len(countArgs)+1)))
		args = append(args, showID)
		countArgs = append(countArgs, showID)
	}

	if isActive != "" {
		active := isActive == "true"
		conditions = append(conditions, "is_active = $"+string(rune(len(args)+1)))
		countConditions = append(countConditions, "is_active = $"+string(rune(len(countArgs)+1)))
		args = append(args, active)
		countArgs = append(countArgs, active)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		baseCountQuery += " AND " + strings.Join(countConditions, " AND ")
	}

	query = baseQuery + " ORDER BY " + orderBy + " " + orderDir + " LIMIT $" + string(rune(len(args)+1)) + " OFFSET $" + string(rune(len(args)+2))
	args = append(args, limit, offset)

	countQuery = baseCountQuery

	pool := db.Pool()
	if pool == nil {
		return utils.NewApiError("Failed to connect to database", 0x0701).AsResponse(c)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("timers query error: %v", err)
		return utils.NewApiError("Failed to query timers", 0x0702).AsResponse(c)
	}
	defer rows.Close()

	timers, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Timer])
	if err != nil {
		log.Printf("timers collection error: %v", err)
		return utils.NewApiError("Failed to process timers", 0x0703).AsResponse(c)
	}

	// Get total count for pagination metadata
	var totalCount int
	err = pool.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		log.Printf("count query error: %v", err)
		totalCount = 0
	}

	totalPages := (totalCount + limit - 1) / limit

	return c.JSON(fiber.Map{
		"timers": timers,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_count": totalCount,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
		"filters": fiber.Map{
			"show_id":   showID,
			"is_active": isActive,
			"order_by":  orderBy,
			"order_dir": orderDir,
		},
	})
}

// GetTimerByID retrieves a single timer by ID
func GetTimerByID(c *fiber.Ctx) error {
	timerID := c.Params("id")
	if timerID == "" {
		return utils.NewApiError("Timer ID is required", 0x0711).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	timer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get timer %s: %v", timerID, err)
		return utils.NewApiError("Timer not found", 0x0712).AsResponse(c)
	}

	return c.JSON(timer)
}
