package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"wanshow-bingo/sse"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/websocket"
)

// aggregateCache stores the latest aggregate JSON payload as raw bytes.
var aggregateCache atomic.Value // holds []byte

// aggHub is the SSE hub used to broadcast aggregate updates to clients.
var aggHub *sse.Hub

func init() {
	aggregateCache.Store([]byte("null"))
}

// GetAggregate serves the latest cached aggregate JSON.
func GetAggregate(c *fiber.Ctx) error {
	b, _ := aggregateCache.Load().([]byte)
	if len(b) == 0 || string(b) == "null" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "aggregate not ready"})
	}
	c.Type("json")
	return c.Send(b)
}

// StartAggregateFetcher now performs a single initial HTTP fetch and then
// maintains a 24/7 websocket connection that updates the aggregate cache.
func StartAggregateFetcher(h *sse.Hub) {
	aggHub = h
	go func() {
		// Initial one-time fetch to warm the cache
		fetchAggregateOnce()
		// Then rely entirely on the websocket stream with reconnection
		connectLiveWebsocketForever()
	}()
}

// fetchAggregateOnce pulls the aggregate endpoint and updates the cache.
func fetchAggregateOnce() {
	base := os.Getenv("WHENPLANE_AGGREGATE_URL")
	if base == "" {
		base = "https://whenplane.com/api/aggregate?fast=true"
	}
	ms := strconv.FormatInt(time.Now().UnixMilli(), 10)
	url := base + "&r=" + ms

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("aggregate fetch error: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("aggregate fetch unexpected status: %d", resp.StatusCode)
		return
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("aggregate read error: %v", err)
		return
	}
	aggregateCache.Store(b)
}

// connectLiveWebsocketForever maintains a perpetual connection to the WhenPlane
// websocket and updates the aggregate cache with any JSON payloads received.
// It ignores ping/pong content and handles reconnection with exponential backoff.
func connectLiveWebsocketForever() {
	backoff := 2 * time.Second
	for {
		wsURL := os.Getenv("WHENPLANE_SOCKET_URL")
		if wsURL == "" {
			log.Println("websocket: WHENPLANE_SOCKET_URL not set; retrying in 30s")
			time.Sleep(30 * time.Second)
			continue
		}

		c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			log.Printf("websocket: dial error: %v", err)
			d := backoff
			if d > 30*time.Second {
				d = 30 * time.Second
			}
			time.Sleep(d)
			if backoff < 30*time.Second {
				backoff *= 2
			}
			continue
		}

		log.Println("websocket: connected")
		backoff = 2 * time.Second

		handleWebsocketMessage(c)

		// small delay before attempting reconnect
		time.Sleep(1 * time.Second)
	}
}

func handleWebsocketMessage(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Printf("websocket: read error: %v", err)
			_ = c.Close()
			break
		}

		// Ignore ping/pong textual keepalives
		if string(message) == "pong" || string(message) == "ping" {
			// Broadcast to hub with opcode whenplane.aggregate if hub is available
			if aggHub != nil {
				evt := sse.SocketEvent{Opcode: "whenplane.aggregate", Data: aggregateCache.Load()}
				if b, err := json.Marshal(evt); err == nil {
					aggHub.Broadcast(string(b))
				}
			}
		}

		// If it looks like JSON, store it directly and broadcast to SSE hub
		if len(message) > 0 && (message[0] == '{' || message[0] == '[') {
			aggregateCache.Store(message)

			log.Println("message received:", string(message))

			// Broadcast to hub with opcode whenplane.aggregate if hub is available
			if aggHub != nil {
				var payload any
				if err := json.Unmarshal(message, &payload); err == nil {
					evt := sse.SocketEvent{Opcode: "whenplane.aggregate", Data: payload}
					if b, err := json.Marshal(evt); err == nil {
						aggHub.Broadcast(string(b))
					}
				}
			}
		}
	}
}
