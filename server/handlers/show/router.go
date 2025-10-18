package show

import (
	"context"
	"wanshow-bingo/db"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/shows", Register)
}

func Register(router fiber.Router) {
	router.Get("/latest", GetLatest)
}

func GetLatest(ctx *fiber.Ctx) error {
	show, err := db.GetLatestShow(context.Background())

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(show)
}
