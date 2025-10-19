package models

import (
	"time"
)

// Player represents a user account
type Player struct {
	ID          string                  `json:"id" db:"id"`
	DID         string                  `json:"did" db:"did"`
	DisplayName string                  `json:"display_name" db:"display_name"`
	Avatar      *string                 `json:"avatar" db:"avatar"`
	Settings    *map[string]interface{} `json:"settings" db:"settings"`
	Score       int                     `json:"score" db:"score"`
	Permissions Permission              `json:"permissions" db:"permissions"`
	CreatedAt   time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time              `json:"deleted_at" db:"deleted_at"`
}

// DiscordUser represents a Discord user from the API
type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	GlobalName    string `json:"global_name"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Verified      bool   `json:"verified"`
}

// Session represents a user session
type Session struct {
	ID        string     `json:"id" db:"id"`
	PlayerID  string     `json:"player_id" db:"player_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ShowState represents the state of a show
type ShowState string

const (
	ShowStateScheduled ShowState = "scheduled"
	ShowStateUpcoming  ShowState = "upcoming"
	ShowStateLive      ShowState = "live"
	ShowStateFinished  ShowState = "finished"
)

// Show represents a WAN show episode
type Show struct {
	ID              string                 `json:"id" db:"id"`
	State           ShowState              `json:"state" db:"state"`
	YoutubeID       *string                `json:"youtube_id" db:"youtube_id"`
	ScheduledTime   *time.Time             `json:"scheduled_time" db:"scheduled_time"`
	ActualStartTime *time.Time             `json:"actual_start_time" db:"actual_start_time"`
	Thumbnail       *string                `json:"thumbnail" db:"thumbnail"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time             `json:"deleted_at" db:"deleted_at"`
}

// Tile represents a bingo tile definition
type Tile struct {
	ID        string                 `json:"id" db:"id"`
	Title     string                 `json:"title" db:"title"`
	Category  *string                `json:"category" db:"category"`
	LastDrawn *time.Time             `json:"last_drawn" db:"last_drawn"`
	CreatedBy *string                `json:"created_by" db:"created_by"`
	Weight    float64                `json:"weight" db:"weight"`
	Score     float64                `json:"score" db:"score"`
	Settings  map[string]interface{} `json:"settings" db:"settings"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at" db:"deleted_at"`
}

// ShowTile represents the junction table linking tiles to shows
type ShowTile struct {
	ShowID    string     `json:"show_id" db:"show_id"`
	TileID    string     `json:"tile_id" db:"tile_id"`
	Weight    float64    `json:"weight" db:"weight"`
	Score     float64    `json:"score" db:"score"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Board represents a player's bingo board for a specific show
type Board struct {
	ID                     string     `json:"id" db:"id"`
	PlayerID               string     `json:"player_id" db:"player_id"`
	ShowID                 string     `json:"show_id" db:"show_id"`
	Tiles                  []string   `json:"tiles" db:"tiles"`
	Winner                 bool       `json:"winner" db:"winner"`
	TotalScore             float64    `json:"total_score" db:"total_score"`
	PotentialScore         float64    `json:"potential_score" db:"potential_score"`
	RegenerationDiminisher float64    `json:"regeneration_diminisher" db:"regeneration_diminisher"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at" db:"deleted_at"`
}

// TileConfirmation records when tiles are confirmed during a show
type TileConfirmation struct {
	ID               string     `json:"id" db:"id"`
	ShowID           string     `json:"show_id" db:"show_id"`
	TileID           string     `json:"tile_id" db:"tile_id"`
	ConfirmedBy      *string    `json:"confirmed_by" db:"confirmed_by"`
	Context          *string    `json:"context" db:"context"`
	ConfirmationTime time.Time  `json:"confirmation_time" db:"confirmation_time"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Message records chat messages during a show
type Message struct {
	ID        string     `json:"id" db:"id"`
	ShowID    string     `json:"show_id" db:"show_id"`
	PlayerID  string     `json:"player_id" db:"player_id"`
	Contents  string     `json:"contents" db:"contents"`
	System    bool       `json:"system" db:"system" `
	Replying  *string    `json:"replying" db:"replying"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Timer represents a countdown timer for shows
type Timer struct {
	ID        string                 `json:"id" db:"id"`
	Title     string                 `json:"title" db:"title"`
	Duration  int                    `json:"duration" db:"duration"`
	CreatedBy *string                `json:"created_by" db:"created_by"`
	ShowID    *string                `json:"show_id" db:"show_id"`
	StartsAt  *time.Time             `json:"starts_at" db:"starts_at"`
	ExpiresAt *time.Time             `json:"expires_at" db:"expires_at"`
	IsActive  bool                   `json:"is_active" db:"is_active"`
	Settings  map[string]interface{} `json:"settings" db:"settings"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at" db:"deleted_at"`
}

// MessageRequest is used in the API to parse successfully
type MessageRequest struct {
	Contents string `json:"contents" db:"contents"`
}

// TileSuggestion represents a user-submitted tile suggestion
type TileSuggestion struct {
	ID         string     `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	TileName   string     `json:"tile_name" db:"tile_name"`
	Reason     string     `json:"reason" db:"reason"`
	Status     string     `json:"status" db:"status"`
	ReviewedBy *string    `json:"reviewed_by" db:"reviewed_by"`
	ReviewedAt *time.Time `json:"reviewed_at" db:"reviewed_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at"`
}
