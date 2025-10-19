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

// CreateTimer creates a new timer
func CreateTimer(c *fiber.Ctx) error {
	// Get authenticated player
	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0721).AsResponse(c)
	}

	var timer models.Timer
	if err := c.BodyParser(&timer); err != nil {
		return utils.NewApiError("Invalid request body", 0x0722).AsResponse(c)
	}

	// Validate required fields
	if timer.Title == "" {
		return utils.NewApiError("Timer title is required", 0x0723).AsResponse(c)
	}
	if timer.Duration <= 0 {
		return utils.NewApiError("Timer duration must be positive", 0x0724).AsResponse(c)
	}

	// Set created_by to current player
	timer.CreatedBy = &player.ID

	// Automatically activate the timer when created
	timer.IsActive = true
	now := time.Now()
	timer.StartsAt = &now
	expiresAt := now.Add(time.Duration(timer.Duration) * time.Second)
	timer.ExpiresAt = &expiresAt

	// If no show_id provided, use latest show
	if timer.ShowID == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		latestShow, err := db.GetLatestShow(ctx)
		if err != nil {
			log.Printf("failed to get latest show: %v", err)
			return utils.NewApiError("Failed to get latest show", 0x0725).AsResponse(c)
		}
		timer.ShowID = &latestShow.ID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.PersistTimer(ctx, &timer)
	if err != nil {
		log.Printf("failed to create timer: %v", err)
		return utils.NewApiError("Failed to create timer", 0x0726).AsResponse(c)
	}

	return c.Status(fiber.StatusCreated).JSON(timer)
}

// UpdateTimer updates an existing timer
func UpdateTimer(c *fiber.Ctx) error {
	timerID := c.Params("id")
	if timerID == "" {
		return utils.NewApiError("Timer ID is required", 0x0731).AsResponse(c)
	}

	// Get authenticated player
	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0732).AsResponse(c)
	}

	var updateData models.Timer
	if err := c.BodyParser(&updateData); err != nil {
		return utils.NewApiError("Invalid request body", 0x0733).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get existing timer
	existingTimer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get timer %s: %v", timerID, err)
		return utils.NewApiError("Timer not found", 0x0734).AsResponse(c)
	}

	// Check if user owns this timer
	if existingTimer.CreatedBy == nil || *existingTimer.CreatedBy != player.ID {
		return utils.NewApiError("Access denied", 0x0735).AsResponse(c)
	}

	// Update fields (only allow certain fields to be updated)
	if updateData.Title != "" {
		existingTimer.Title = updateData.Title
	}
	if updateData.Duration > 0 {
		existingTimer.Duration = updateData.Duration
	}
	if updateData.Settings != nil {
		existingTimer.Settings = updateData.Settings
	}

	err = db.PersistTimer(ctx, existingTimer)
	if err != nil {
		log.Printf("failed to update timer: %v", err)
		return utils.NewApiError("Failed to update timer", 0x0736).AsResponse(c)
	}

	return c.JSON(existingTimer)
}

// DeleteTimer soft deletes a timer
func DeleteTimer(c *fiber.Ctx) error {
	timerID := c.Params("id")
	if timerID == "" {
		return utils.NewApiError("Timer ID is required", 0x0741).AsResponse(c)
	}

	// Get authenticated player
	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0742).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get existing timer to check ownership
	existingTimer, err := db.GetTimerByID(ctx, timerID)
	if err != nil {
		log.Printf("failed to get timer %s: %v", timerID, err)
		return utils.NewApiError("Timer not found", 0x0743).AsResponse(c)
	}

	// Check if user owns this timer
	if existingTimer.CreatedBy == nil || *existingTimer.CreatedBy != player.ID {
		return utils.NewApiError("Access denied", 0x0744).AsResponse(c)
	}

	err = db.DeleteTimer(ctx, timerID)
	if err != nil {
		log.Printf("failed to delete timer: %v", err)
		return utils.NewApiError("Failed to delete timer", 0x0745).AsResponse(c)
	}

	return c.JSON(fiber.Map{"success": true})
}
