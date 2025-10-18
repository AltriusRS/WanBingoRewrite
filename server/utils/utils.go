package utils

import (
	"log"
	"wanshow-bingo/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var app *fiber.App

func init() {
	app = fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://app.bingo.local,https://api.bingo.local,https://discord.com,http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000,http://127.0.0.1:3001", // your frontend origin, API origin, Discord, plus common dev ports
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // REQUIRED for cookies or Auth headers
	}))

	app.Use(middleware.RequestLogger)

	// Whoami endpoint to check current user
	app.Get("/whoami", middleware.OptionalPlayerAuthMiddleware, func(c *fiber.Ctx) error {
		player := c.Locals("player")
		if player != nil {
			return c.JSON(fiber.Map{"user": player})
		} else {
			return c.JSON(fiber.Map{"user": nil})
		}
	})

}

func RegisterRouter(path string, callback func(c fiber.Router)) {
	router := app.Group(path)
	callback(router)
}

func GetApp() *fiber.App {
	return app
}

func StartRouter(port string) {
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

type ApiError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewApiError(message string, code int) *ApiError {
	return &ApiError{
		Message: message,
		Code:    code,
	}
}

func (e *ApiError) AsResponse(c *fiber.Ctx) error {
	return c.Status(500).JSON(e)
}
