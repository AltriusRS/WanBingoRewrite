package stream

import (
	"wanshow-bingo/sse"

	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	client := sse.NewClient()
	client.Hub = sse.GetHostHub()
	return client.Bind(c)
}
