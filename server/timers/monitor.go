package timers

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/sse"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/robfig/cron/v3"
)

var timerCron *cron.Cron

func init() {
	timerCron = cron.New()
	timerCron.Start()

	// Check for expired timers every 2 seconds
	_, err := timerCron.AddFunc("@every 2s", checkExpiredTimers)
	if err != nil {
		log.Printf("Failed to schedule timer monitor: %v", err)
	}

	log.Println("Timer monitor initialized")
}

// confirmWANTile automatically confirms the 4-hour WAN show tile
func confirmWANTile(ctx context.Context, showID string) {
	tileID := "BfaqFYztlR"
	contextStr := "Automatic confirmation after 4-hour timer"

	id, _ := gonanoid.New(10)

	confirmation := &models.TileConfirmation{
		ID:               id,
		ShowID:           showID,
		TileID:           tileID,
		ConfirmedBy:      nil, // system
		Context:          &contextStr,
		ConfirmationTime: time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := db.PersistTileConfirmation(ctx, confirmation)
	if err != nil {
		log.Printf("Failed to confirm WAN tile: %v", err)
		return
	}

	// Create system message
	messageContent := "**TILE CONFIRMED** 4 Hour WAN Show"
	systemMessage := &models.Message{
		ID:        uuid.New().String(),
		ShowID:    showID,
		PlayerID:  "SYSTEM",
		Contents:  messageContent,
		System:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.PersistMessage(ctx, systemMessage)
	if err != nil {
		log.Printf("Failed to send WAN tile confirmation message: %v", err)
		return
	}

	// Broadcast message to chat hub
	chatHub := sse.GetChatHub()
	if chatHub != nil {
		chatHub.BroadcastEvent("chat.message", systemMessage)
	}

	// Broadcast tile confirmation to host hub
	hostHub := sse.GetHostHub()
	if hostHub != nil {
		hostHub.BroadcastEvent("tile.confirm", map[string]interface{}{
			"tileId": tileID,
		})
	}

	log.Printf("Automatically confirmed WAN tile for show %s", showID)
}

// checkExpiredTimers checks for expired timers and sends SSE events
func checkExpiredTimers() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	expiredTimers, err := db.GetExpiredTimers(ctx)
	if err != nil {
		log.Printf("Failed to get expired timers: %v", err)
		return
	}

	for _, timer := range expiredTimers {
		// Special handling for WAN Show Timer
		if timer.Title == "WAN Show Timer" && timer.ShowID != nil {
			// Check if show is still live
			latestShow, err := db.GetLatestShow(ctx)
			if err == nil && latestShow.State == models.ShowStateLive && latestShow.ID == *timer.ShowID {
				// Automatically confirm the 4-hour WAN show tile
				confirmWANTile(ctx, latestShow.ID)
			}
		}

		// Send timer.expired event to the host SSE stream
		hostHub := sse.GetHostHub()
		if hostHub != nil {
			hostHub.BroadcastEvent("timer.expired", map[string]interface{}{
				"timer_id":   timer.ID,
				"title":      timer.Title,
				"show_id":    timer.ShowID,
				"created_by": timer.CreatedBy,
				"expired_at": timer.ExpiresAt,
			})
		}

		// Stop the timer (deactivate it)
		err := db.StopTimer(ctx, timer.ID)
		if err != nil {
			log.Printf("Failed to stop expired timer %s: %v", timer.ID, err)
		}
	}
}

// Cleanup stops the timer monitor
func Cleanup() {
	if timerCron != nil {
		ctx := timerCron.Stop()
		<-ctx.Done()
	}
}
