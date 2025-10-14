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
		AllowOrigins:     "http://localhost:7999", // your frontend origin
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // REQUIRED for cookies or Auth headers
	}))

	app.Use(middleware.RequestLogger)

	// Example protected route
	app.Get("/whoami", middleware.RequiredAuthMiddleware, func(c *fiber.Ctx) error {
		claims := c.Locals("claims")
		return c.JSON(claims)
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

	routes := app.GetRoutes()

	log.Println("Registered routes:")
	for _, route := range routes {
		log.Println(route.Method, route.Path)
	}
	log.Println("")

	log.Printf("🚀 Server running on http://localhost:%s", port)
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
