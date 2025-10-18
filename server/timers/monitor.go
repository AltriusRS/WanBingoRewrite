package timers

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/sse"

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
