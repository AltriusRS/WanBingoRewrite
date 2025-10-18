package socket

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"wanshow-bingo/whenplane"
	"wanshow-bingo/whenplane/watcher"

	"github.com/gorilla/websocket"
)

func Init() {
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

	log.Printf("Fetching aggregate from %s", url)

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

	log.Printf("aggregate fetched\n")

	aggregate, err := whenplane.AggregateFromJSON(string(b))

	if err != nil {
		log.Println("error parsing aggregate response:", err)
		log.Println("body:", string(b))
		return
	}

	watcher.AggregateChan <- &aggregate
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
	log.Println("websocket: handleWebsocketMessage")
	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Printf("websocket: read error: %v", err)
			_ = c.Close()
			break
		}

		// Ignore ping/pong textual keepalives
		if string(message) == "pong" || string(message) == "ping" {
			continue
		}

		// If it looks like JSON, store it directly and broadcast to SSE hub
		if len(message) > 0 && (message[0] == '{' || message[0] == '[') {
			AggregateMessage(string(message))
		}
	}
}

func AggregateMessage(message string) {
	aggregate, err := whenplane.AggregateFromJSON(message)

	if err != nil {
		log.Println("error parsing aggregate event:", err)
		log.Println("message received:", message)
		return
	}

	watcher.AggregateChan <- &aggregate
}
