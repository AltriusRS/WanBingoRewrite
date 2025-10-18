package aggregaterouter

import (
	"wanshow-bingo/whenplane"

	"github.com/gofiber/fiber/v2"
)

// Get serves the latest cached aggregate JSON.
func Get(c *fiber.Ctx) error {
	aggregate, err := whenplane.GetAggregateCache()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
	return c.JSON(aggregate)
}
