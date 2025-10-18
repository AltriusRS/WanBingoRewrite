package sse

import (
	"bufio"
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/whenplane"

	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Client struct {
	Id              string
	Hub             *Hub
	Queue           chan string
	Writer          *bufio.Writer
	ticker          *time.Ticker
	IsAuthenticated bool
	Player          *models.Player
}

func NewClient() *Client {
	id, err := gonanoid.New(10)

	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Id:              id,
		Queue:           make(chan string, 10),
		IsAuthenticated: false,
	}
}

func (c *Client) Bind(ctx *fiber.Ctx) error {
	playerRaw := ctx.Locals("player")

	if playerRaw != nil {
		c.IsAuthenticated = true
		if player, ok := playerRaw.(*models.Player); ok {
			c.Player = player
		}
	}

	// Set required headers for SSE on the RESPONSE
	ctx.Set("Content-Type", "text/event-stream; charset=utf-8")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("X-Accel-Buffering", "no") // Disable proxy buffering if applicable

	// Set the streaming body Writer which will handle the entire SSE lifecycle.
	ctx.Context().SetBodyStreamWriter(c.streamWriter)
	return nil
}

// streamWriter runs for the lifetime of the response. It should not return
// until the client disconnects or an error occurs.
func (c *Client) streamWriter(w *bufio.Writer) {
	c.Writer = w

	// Register client and start keep-alive ticker
	c.ticker = time.NewTicker(30 * time.Second)
	if c.Hub != nil {
		c.Hub.RegisterClient(c)
	}
	defer c.Stop()

	// Send the initial connected event
	connectEvent := BuildEvent("hub.connected", "Connected to chat hub.")
	go c.Send(connectEvent.String())

	EmitWhenplaneAgg(c)

	if c.Hub != nil {
		authenticatedEvent := BuildEvent("hub.authenticated", ClientCapabilities{
			CanChat: true,
		})
		go c.Send(authenticatedEvent.String())

		if c.Hub.name == "CHAT" {
			go SendChatHistory(c)
		}
	}

	// Listen for messages and keep-alive ticks
	for {
		select {
		case message, ok := <-c.Queue:
			if !ok {
				return
			}
			if err := c.write(message); err != nil {
				return
			}
		case <-c.ticker.C:
			var msg string
			if c.Hub != nil {
				msg = c.Hub.BuildConnectionCount()
			}
			if msg != "" {
				if err := c.write(msg); err != nil {
					return
				}
			}
		}
	}
}

func (c *Client) Send(message string) {
	c.Queue <- message
}

func (c *Client) Read() <-chan string {
	return c.Queue
}

func (c *Client) Stop() {
	if c.Hub != nil {
		c.Hub.UnregisterClient(c)
	}
	if c.ticker != nil {
		c.ticker.Stop()
	}
}

// write sends a single SSE data frame and flushes the Writer.
func (c *Client) write(msg string) error {
	if c.Writer == nil {
		return nil
	}
	_, err := c.Writer.WriteString("data: " + msg + "\n\n")
	if err != nil {
		log.Println("[SSE ClientChannel] - Error writing to Writer:", err)
		return err
	}
	if err := c.Writer.Flush(); err != nil {
		log.Println("[SSE ClientChannel] - Error flushing Writer:", err)
		return err
	}

	return nil
}

func EmitWhenplaneAgg(c *Client) {
	agg, err := whenplane.GetAggregateCache()

	if err == nil {
		aggregateEvent := BuildEvent("whenplane.aggregate", agg)
		go c.Send(aggregateEvent.String())
	} else {
		log.Println("[SSE ClientChannel] - Error getting aggregate cache:", err)
	}
}

type ClientCapabilities struct {
	CanChat     bool `json:"canChat"`
	CanHost     bool `json:"canHost"`
	CanModerate bool `json:"canModerate"`
}

var UnauthorizedCapabilities = ClientCapabilities{
	CanChat:     false,
	CanHost:     false,
	CanModerate: false,
}

func (c *Client) GetCapabilities() ClientCapabilities {
	if !c.IsAuthenticated {
		return UnauthorizedCapabilities
	}

	pool := db.Pool()

	if pool == nil {
		return UnauthorizedCapabilities
	}

	return UnauthorizedCapabilities
}

func SendChatHistory(c *Client) {
	history, err := db.GetMessageHistory(context.Background())

	if err != nil {
		log.Printf("[SSE ClientChannel] - Failed to retrieve chat history - %v", err)
		return
	}

	// Get online players from the hub
	playerMap := make(map[string]map[string]interface{})
	if c.Hub != nil {
		for _, client := range c.Hub.GetClients() {
			if client.IsAuthenticated && client.Player != nil {
				player := client.Player
				playerMap[player.ID] = map[string]interface{}{
					"id":           player.ID,
					"display_name": player.DisplayName,
					"avatar":       player.Avatar,
					"permissions":  player.Permissions,
					"settings":     player.Settings,
					"created_at":   player.CreatedAt,
				}
			}
		}
	}

	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	// Send chat history
	for _, msg := range history {
		// Attach player info to message
		messageWithPlayer := map[string]interface{}{
			"id":         msg.ID,
			"show_id":    msg.ShowID,
			"player_id":  msg.PlayerID,
			"contents":   msg.Contents,
			"system":     msg.System,
			"replying":   msg.Replying,
			"created_at": msg.CreatedAt,
			"updated_at": msg.UpdatedAt,
			"deleted_at": msg.DeletedAt,
			"player":     playerMap[msg.PlayerID],
		}
		msgEvent := BuildEvent("chat.message", messageWithPlayer)
		c.Send(msgEvent.String())
	}

	// Send player information for the chat participants
	playerList := make([]map[string]interface{}, 0, len(playerMap))
	for _, playerInfo := range playerMap {
		playerList = append(playerList, playerInfo)
	}

	playersEvent := BuildEvent("chat.players", map[string]interface{}{
		"players": playerList,
		"count":   len(playerList),
	})
	c.Send(playersEvent.String())
}
