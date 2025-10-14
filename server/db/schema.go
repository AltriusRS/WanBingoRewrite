package db

import (
	"time"
)

type Tile struct {
	Id            string    `json:"id"`
	Text          string    `json:"text"`            // The text of the tile to be displayed
	Category      string    `json:"category"`        // The category of the tile
	Weight        float64   `json:"weight"`          // The weight of the tile - Higher weight means more likely to be in a board
	LastDrawnShow string    `json:"last_drawn_show"` // The last show ID that this tile was drawn in
	IsActive      bool      `json:"is_active"`       // Whether the tile is active or not
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Show struct {
	Id           string    `json:"id"`
	YoutubeId    string    `json:"yt_id"`        // The YouTube ID of the show
	FloatplaneId string    `json:"fp_id"`        // The Floatplane ID of the show (if applicable)
	StartsAt     time.Time `json:"starts_at"`    // The scheduled start time of the show
	WentLiveAt   time.Time `json:"went_live_at"` // The actual time the show went live
	Title        string    `json:"title"`        // The title of the show (with date mark removed)
	IsLive       bool      `json:"is_live"`      // Whether the show is currently live
	Metadata     any       `json:"metadata"`     // Additional metadata about the show (free-structured JSON
	CreatedAt    time.Time `json:"created_at"`
}

type TileConfirmation struct {
	Id          string    `json:"id"`
	ShowId      string    `json:"show_id"`      // The show ID that the tile was drawn in
	TileId      string    `json:"tile_id"`      // The tile ID that was drawn
	Context     string    `json:"context"`      // The context provided by the host that confirmed the tile
	ConfirmedBy string    `json:"confirmed_by"` // The user that confirmed the tile
	ConfirmedAt time.Time `json:"confirmed_at"` // The time the tile was confirmed
}

type HostLock struct {
	TileId    string    `json:"tile_id"`    // The tile ID that was locked
	ShowId    string    `json:"show_id"`    // The show ID that the tile was locked in
	LockedBy  string    `json:"locked_by"`  // The user that locked the tile
	ExpiresAt time.Time `json:"expires_at"` // The time the lock expires
	CreatedAt time.Time `json:"created_at"`
}

type ChatMessage struct {
	Id        string    `json:"id"`
	ShowId    string    `json:"show_id"`  // The show ID that the message was sent in
	Type      string    `json:"type"`     // "user" | "system"
	Username  string    `json:"username"` // Username of the sender, if applicable
	Message   string    `json:"message"`  // Message content
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tile) TableName() string {
	return "tile_router"
}

func (s *Show) TableName() string {
	return "shows"
}

func (c *ChatMessage) TableName() string {
	return "chat_messages"
}

func (l *HostLock) TableName() string {
	return "host_locks"
}

func (c *TileConfirmation) TableName() string {
	return "tile_confirmations"
}
