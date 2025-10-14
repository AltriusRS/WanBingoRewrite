package main

import (
	"os"
	"wanshow-bingo/db"
	_ "wanshow-bingo/handlers"
	_ "wanshow-bingo/sse"
	"wanshow-bingo/utils"

	"github.com/workos/workos-go/pkg/audittrail"
	"github.com/workos/workos-go/pkg/directorysync"
	"github.com/workos/workos-go/pkg/organizations"
	"github.com/workos/workos-go/pkg/passwordless"
	"github.com/workos/workos-go/pkg/portal"
	"github.com/workos/workos-go/pkg/sso"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
)

func init() {
	// Configure WorkOS client.
	wosApiKey := os.Getenv("WORKOS_API_KEY")
	wosClientId := os.Getenv("WORKOS_CLIENT_ID")

	sso.Configure(wosApiKey, wosClientId)
	organizations.SetAPIKey(wosApiKey)
	passwordless.SetAPIKey(wosApiKey)
	directorysync.SetAPIKey(wosApiKey)
	audittrail.SetAPIKey(wosApiKey)
	portal.SetAPIKey(wosApiKey)
	usermanagement.SetAPIKey(wosApiKey)
}

func main() {

	// Initialize optional database pool.
	db.Init()

	//app.Post("/chat/system", func(c *fiber.Ctx) error {
	//	password := os.Getenv("HOST_PASSWORD")
	//	if c.Get("Authorization") != "Bearer "+password {
	//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	//	}
	//	return handlers.SendSystemMessage(c, hub)
	//})
	//
	//app.Post("/chat/message", func(c *fiber.Ctx) error {
	//	// TODO: Make this endpoint require authorization of some variety.
	//	return handlers.SendChatMessage(c, hub)
	//})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.StartRouter(port)

}
