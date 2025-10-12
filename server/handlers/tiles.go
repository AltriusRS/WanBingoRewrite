package handlers

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"wanshow-bingo/db"

	"github.com/gofiber/fiber/v2"
)

type Tile struct {
	Text          string  `json:"text"`
	Category      string  `json:"category"`
	Weight        float64 `json:"weight"`
	LastDrawnShow int     `json:"last_drawn_show"`
}

func GetTiles(c *fiber.Ctx) error {
	// Try database first if available.
	if pool := db.Pool(); pool != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		rows, err := pool.Query(ctx, `
			SELECT text, category, weight, COALESCE(last_drawn_show, 0)
			FROM tiles
			ORDER BY random()
			LIMIT 25`)
		if err == nil {
			defer rows.Close()
			var tiles []Tile
			for rows.Next() {
				var t Tile
				if err := rows.Scan(&t.Text, &t.Category, &t.Weight, &t.LastDrawnShow); err != nil {
					log.Printf("tiles query scan error: %v", err)
					tiles = nil
					break
				}
				tiles = append(tiles, t)
			}
			if tiles != nil && len(tiles) > 0 {
				return c.JSON(tiles)
			}
		} else {
			log.Printf("tiles query error (falling back to file): %v", err)
		}
	}

	// Fallback to bundled JSON file.
	file, err := os.ReadFile("data/tiles.json")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "tiles not found"})
	}

	var allTiles []Tile
	if err := json.Unmarshal(file, &allTiles); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "invalid tiles data"})
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allTiles), func(i, j int) { allTiles[i], allTiles[j] = allTiles[j], allTiles[i] })

	n := 25
	if len(allTiles) < n {
		n = len(allTiles)
	}

	return c.JSON(allTiles[:n])
}
