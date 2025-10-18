package main

import (
	"os"
	_ "time/tzdata"
	"wanshow-bingo/db"
	_ "wanshow-bingo/handlers"
	"wanshow-bingo/middleware"
	_ "wanshow-bingo/sse"
	_ "wanshow-bingo/timers"
	"wanshow-bingo/utils"
	"wanshow-bingo/whenplane/socket"
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.StartRouter(port)

}
