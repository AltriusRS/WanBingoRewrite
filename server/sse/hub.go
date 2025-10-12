package sse

import "encoding/json"

// Simple SSE hub for broadcasting JSON strings to connected clients.
// The handlers expect to send/receive raw strings already formatted as SSE data lines.

// Client represents a connection-specific channel for outgoing messages.
type Client chan string

type SocketEvent struct {
	Opcode string `json:"opcode"`
	Data   any    `json:"data"`
}

type MemberCount struct {
	Count int `json:"count"`
}

type Hub struct {
	clients    map[Client]bool
	register   chan Client
	unregister chan Client
	broadcast  chan string
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[Client]bool),
		register:   make(chan Client),
		unregister: make(chan Client),
		broadcast:  make(chan string),
	}
}

// Run processes register/unregister/broadcast events.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c)
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				// Format as SSE event data lines the writer expects
				select {
				case c <- "data: " + msg + "\n\n":
				default:
					// client is blocked or gone; clean up
					delete(h.clients, c)
					close(c)
				}
			}
		}
	}
}
func (h *Hub) RegisterClient(c Client) {
	h.BroadcastMemberCount()
	h.register <- c
}

func (h *Hub) UnregisterClient(c Client) {
	h.BroadcastMemberCount()
	h.unregister <- c
}

func (h *Hub) Broadcast(msg string) {
	h.broadcast <- msg
}

func (h *Hub) BroadcastMemberCount() {
	h.Broadcast(h.BuildMemberCount())
}

func (h *Hub) BuildMemberCount() string {
	var count_event SocketEvent
	count_event.Opcode = "chat.members.count"
	count_event.Data = MemberCount{Count: len(h.clients)}

	b, _ := json.Marshal(count_event)

	return string(b)
}
