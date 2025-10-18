package aggregaterouter

import (
	"os"
	"wanshow-bingo/whenplane"
	"wanshow-bingo/whenplane/watcher"

	"github.com/gofiber/fiber/v2"
)

func Post(c *fiber.Ctx) error {
	password := os.Getenv("HOST_PASSWORD")
	if c.Get("Authorization") != "Bearer "+password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var payload whenplane.Aggregate

	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	watcher.AggregateChan <- &payload
	return c.SendStatus(fiber.StatusOK)
}
