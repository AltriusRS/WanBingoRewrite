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
	router.Get("/:id", GetByID)
}

func GetLatest(ctx *fiber.Ctx) error {
	show, err := db.GetLatestShow(context.Background())

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(show)
}

func GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Show ID is required")
	}

	show, err := db.GetShowByID(context.Background(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Show not found")
	}

	return ctx.JSON(show)
}
