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

// TODO: Add filters + pagination support

func Get(c *fiber.Ctx) error {
	pool := db.Pool()

	// If no database is available, return an error
	if pool == nil {
		return utils.NewApiError("Failed to connect to database", 0x0101).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := pool.Query(ctx, `
		SELECT *
		FROM tiles
		ORDER BY random()
	`)

	if err != nil {
		log.Printf("tiles query error: %v", err)
		return utils.NewApiError("Failed to query tiles", 0x0201).AsResponse(c)
	}

	defer rows.Close()
	tiles, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Tile])

	if err != nil {
		log.Printf("tiles query error: %v", err)
	}

	return c.JSON(tiles)
}
