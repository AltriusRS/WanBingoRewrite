package aggregaterouter

import (
	"wanshow-bingo/whenplane"

	"github.com/gofiber/fiber/v2"
)

// Get serves the latest cached aggregate JSON.
func Get(c *fiber.Ctx) error {
	aggregate := whenplane.GetAggregateCache()
	return c.JSON(aggregate)
}
