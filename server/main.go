package main

import (
	"os"
	"wanshow-bingo/db"
	_ "wanshow-bingo/handlers"
	"wanshow-bingo/middleware"
	_ "wanshow-bingo/sse"
	"wanshow-bingo/utils"
	"wanshow-bingo/whenplane/socket"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func init() {
	// Initialize optional database pool.
	db.Init()

	// Initialize the whenplane socket aggregator.
	socket.Init()

	// Initialize Discord OAuth
	middleware.InitDiscordOAuth()
}

func main() {

	// Register auth routes
	utils.RegisterRouter("/auth", func(router fiber.Router) {

		// Protected Discord routes
		discord := router.Group("/discord")
		discord.Use(middleware.DiscordAuthMiddleware)
		discord.Get("/user", func(c *fiber.Ctx) error {
			user, err := middleware.GetDiscordUserFromContext(c)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(utils.NewApiError("Not authenticated with Discord", 401))
			}

			return c.JSON(fiber.Map{
				"success": true,
				"user": fiber.Map{
					"id":            user.ID,
					"username":      user.Username,
					"discriminator": user.Discriminator,
					"email":         user.Email,
					"avatar":        user.Avatar,
					"verified":      user.Verified,
				},
			})
		})

		discord.Get("/guilds", func(c *fiber.Ctx) error {
			discordToken := c.Cookies("discord-token")
			if discordToken == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(utils.NewApiError("Not authenticated with Discord", 401))
			}

			// Parse the token from the cookie
			token := &oauth2.Token{
				AccessToken: discordToken,
				TokenType:   "Bearer",
			}

			// Get Discord guilds
			guilds, err := middleware.GetDiscordGuilds(token)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to fetch Discord guilds", 500))
			}

			return c.JSON(fiber.Map{
				"success": true,
				"guilds":  guilds,
			})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.StartRouter(port)

}
