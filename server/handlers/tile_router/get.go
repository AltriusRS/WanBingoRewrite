package tilerouter

import (
	"context"
	"log"
	"time"

	"wanshow-bingo/db"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// TODO: Add filters + pagination support

func Get(c *fiber.Ctx) error {
	//filter, err := NewTileFilterFromCtx(c)
	//
	//if err != nil {
	//	return err
	//}
	//
	//return c.JSON(filter)

	pool := db.Pool()

	// If no database is available, return an error
	if pool == nil {
		return utils.NewApiError("Failed to connect to database", 0x0101).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := pool.Query(ctx, `
		SELECT id, text, category, weight, last_drawn_show, is_active
		FROM tiles
		WHERE is_active = true
		ORDER BY weight DESC
	`)

	if err != nil {
		log.Printf("tile_router query error: %v", err)
		return utils.NewApiError("Failed to query tile_router", 0x0201).AsResponse(c)
	}

	defer rows.Close()
	var tiles []db.Tile
	for rows.Next() {
		var t db.Tile
		if err := rows.Scan(&t.Id, &t.Text, &t.Category, &t.Weight, &t.LastDrawnShow, &t.IsActive); err != nil {
			log.Printf("tile_router query scan error: %v", err)
			tiles = nil
			break
		}
		tiles = append(tiles, t)
	}

	return c.JSON(tiles)
}

type TileFilter struct {
	Category      string `json:"category_gte"`
	Id            string `json:"id"`
	Text          string `json:"text"`
	WeightGte     string `json:"weight_gte"`
	WeightLte     string `json:"weight_lte"`
	Weight        string `json:"weight"`
	LastDrawnShow string `json:"last_drawn_show"`
	IsActive      bool   `json:"is_active" default:"true"`
	Limit         int    `json:"limit" default:"100"`
	Offset        int    `json:"offset" default:"0"`
	Sort          string `json:"sort" default:"weight_descending"`
	Mode          string `json:"mode" default:"sheet"` // "host" | "sheet"
}

func NewTileFilterFromCtx(c *fiber.Ctx) (*TileFilter, error) {
	filter := &TileFilter{}

	err := c.QueryParser(filter)

	if err != nil {
		return nil, err
	}

	log.Println(filter)

	return filter, nil
}
