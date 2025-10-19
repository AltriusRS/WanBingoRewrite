package timers

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// StartTimer activates a timer
func StartTimer(c *fiber.Ctx) error {
	timerID := c.Params("id")
	if timerID == "" {
		return utils.NewApiError("Timer ID is required", 0x0751).AsResponse(c)
	}

	// Get authenticated player
	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0752).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get existing timer to check ownership
	existingTimer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get timer %s: %v", timerID, err)
		return utils.NewApiError("Timer not found", 0x0753).AsResponse(c)
	}

	// Check if user owns this timer
	if existingTimer.CreatedBy == nil || *existingTimer.CreatedBy != player.ID {
		return utils.NewApiError("Access denied", 0x0754).AsResponse(c)
	}

	err = db.StartTimer(ctx, timerID)
	if err != nil {
		log.Printf("failed to start timer: %v", err)
		return utils.NewApiError("Failed to start timer", 0x0755).AsResponse(c)
	}

	// Get updated timer
	updatedTimer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get updated timer: %v", err)
		return utils.NewApiError("Timer started but failed to retrieve", 0x0756).AsResponse(c)
	}

	return c.JSON(updatedTimer)
}

// StopTimer deactivates a timer
func StopTimer(c *fiber.Ctx) error {
	timerID := c.Params("id")
	if timerID == "" {
		return utils.NewApiError("Timer ID is required", 0x0761).AsResponse(c)
	}

	// Get authenticated player
	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0762).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get existing timer to check ownership
	existingTimer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get timer %s: %v", timerID, err)
		return utils.NewApiError("Timer not found", 0x0763).AsResponse(c)
	}

	// Check if user owns this timer
	if existingTimer.CreatedBy == nil || *existingTimer.CreatedBy != player.ID {
		return utils.NewApiError("Access denied", 0x0764).AsResponse(c)
	}

	err = db.StopTimer(ctx, timerID)
	if err != nil {
		log.Printf("failed to stop timer: %v", err)
		return utils.NewApiError("Failed to stop timer", 0x0765).AsResponse(c)
	}

	// Get updated timer
	updatedTimer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get updated timer: %v", err)
		return utils.NewApiError("Timer stopped but failed to retrieve", 0x0766).AsResponse(c)
	}

	return c.JSON(updatedTimer)
}

// ResetTimer stops and restarts a timer
func ResetTimer(c *fiber.Ctx) error {
	timerID := c.Params("id")
	if timerID == "" {
		return utils.NewApiError("Timer ID is required", 0x0771).AsResponse(c)
	}

	// Get authenticated player
	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0772).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get existing timer to check ownership
	existingTimer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get timer %s: %v", timerID, err)
		return utils.NewApiError("Timer not found", 0x0773).AsResponse(c)
	}

	// Check if user owns this timer
	if existingTimer.CreatedBy == nil || *existingTimer.CreatedBy != player.ID {
		return utils.NewApiError("Access denied", 0x0774).AsResponse(c)
	}

	// Reset by starting the timer again (this will recalculate start/expire times)
	err = db.StartTimer(ctx, timerID)
	if err != nil {
		log.Printf("failed to reset timer: %v", err)
		return utils.NewApiError("Failed to reset timer", 0x0775).AsResponse(c)
	}

	// Get updated timer
	updatedTimer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get updated timer: %v", err)
		return utils.NewApiError("Timer reset but failed to retrieve", 0x0776).AsResponse(c)
	}

	return c.JSON(updatedTimer)
}
