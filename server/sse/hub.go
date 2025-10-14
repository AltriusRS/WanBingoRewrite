package sse

import (
	"encoding/json"
	"log"
	"wanshow-bingo/utils"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Simple SSE hub for broadcasting JSON strings to connected clients.
// The handlers expect to send/receive raw strings already formatted as SSE data lines.

type SocketEvent struct {
	Id     string `json:"id"`
	Opcode string `json:"opcode"`
	Data   any    `json:"data"`
}

func BuildEvent(opcode string, data any) SocketEvent {
	id, err := gonanoid.New(10)

	if err != nil {
		log.Fatal(err)
	}

	return SocketEvent{
		Id:     id,
		Opcode: opcode,
		Data:   data,
	}
}

func (e *SocketEvent) String() string {
	b, _ := e.MarshalJSON()
	return string(b)
}

func (e *SocketEvent) MarshalJSON() ([]byte, error) {
	// Avoid infinite recursion: json.Marshal will call MarshalJSON again on e.
	// Use an alias type without the method set.
	type alias SocketEvent
	return json.Marshal((*alias)(e))
}

func (e *SocketEvent) Broadcast(hub *Hub) error {
	utils.Debugf("[SSE - %s] Broadcasting event (0x01): %+v", hub.name, e)
	hub.Broadcast(e.String())
	return nil
}

type MemberCount struct {
	Count int `json:"count"`
}

type Hub struct {
	name       string
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan string
}

func NewHub(name string) *Hub {
	return &Hub{
		name:       name,
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan string, 256),
	}
}

// Run processes register/unregister/broadcast events.
func (h *Hub) Run() {
	utils.Debugf("[SSE - %s] Starting hub handler", h.name)
	for {
		select {
		case c := <-h.register:
			utils.Debugf("[SSE - %s] Registering client %s", h.name, c)
			h.clients[c.Id] = c
			go h.BroadcastConnectionCount()
		case c := <-h.unregister:
			if _, ok := h.clients[c.Id]; ok {
				utils.Debugf("[SSE - %s] Deregistering client %s", h.name, c)
				delete(h.clients, c.Id)
				go h.BroadcastConnectionCount()
			}
		case msg := <-h.broadcast:
			utils.Debugf("[SSE - %s] (0x03) Broadcasting message to %d clients: %s", h.name, len(h.clients), msg)
			for c := range h.clients {
				client := h.clients[c]

				// Format as SSE event data lines the Writer expects
				select {
				case client.Queue <- msg:
					utils.Debugf("[SSE - %s] (0x04) Sent message to client %s", h.name, c)
				default:
					// Inline unregister to avoid deadlock from sending to unregister channel within Run
					if _, ok := h.clients[c]; ok {
						utils.Debugf("[SSE - %s] ClientChannel %s not ready; unregistering", h.name, c)
						go h.UnregisterClient(client)
					}
				}
			}
		}
	}
}
func (h *Hub) RegisterClient(c *Client) {
	h.register <- c
}

func (h *Hub) UnregisterClient(c *Client) {
	h.unregister <- c
}

func (h *Hub) Broadcast(msg string) {
	utils.Debugf("[SSE - %s] Broadcasting message (0x02): %s", h.name, msg)
	h.broadcast <- msg
}

func (h *Hub) BroadcastConnectionCount() {
	h.Broadcast(h.BuildConnectionCount())
}

func (h *Hub) GetClient(id string) *Client {
	return h.clients[id]
}

func (h *Hub) GetClients() map[string]*Client {
	return h.clients
}

func (h *Hub) GetConnectedUserList() []string {
	list := make([]string, 0, len(h.clients))
	for id := range h.clients {
		list = append(list, id)
	}
	return list
}

func (h *Hub) BuildConnectionCount() string {
	countEvent := BuildEvent("hub.connections.count", MemberCount{Count: len(h.clients)})
	return countEvent.String()
}

func (h *Hub) BroadcastEvent(eventName string, data any) {
	event := BuildEvent(eventName, data)
	utils.Debugf("[SSE - %s] Building Broadcast Event: %+v", h.name, event)
	err := event.Broadcast(h)
	if err != nil {
		log.Printf("[SSE - %s] Broadcast error: %v - %+v", h.name, err, event)
	}
}
