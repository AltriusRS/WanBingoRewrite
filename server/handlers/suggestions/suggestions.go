package suggestions

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/suggestions", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	router.Post("/", CreateSuggestion)
	router.Get("/", middleware.AuthMiddleware, GetSuggestions)
	router.Put("/:id", middleware.AuthMiddleware, UpdateSuggestion)
}

// CreateSuggestion handles POST /api/suggestions
func CreateSuggestion(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req struct {
		Name     string `json:"name"`
		TileName string `json:"tileName"`
		Reason   string `json:"reason"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.NewApiError("Invalid request body", 0x0601).AsResponse(c)
	}

	if req.Name == "" || req.TileName == "" || req.Reason == "" {
		return utils.NewApiError("Name, tile name, and reason are required", 0x0602).AsResponse(c)
	}

	suggestion, err := db.CreateTileSuggestion(ctx, req.Name, req.TileName, req.Reason)
	if err != nil {
		log.Printf("failed to create suggestion: %v", err)
		return utils.NewApiError("Failed to create suggestion", 0x0603).AsResponse(c)
	}

	return c.Status(201).JSON(suggestion)
}

// GetSuggestions handles GET /api/suggestions
func GetSuggestions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get authenticated player (host only?)
	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0401).AsResponse(c)
	}

	// TODO: Check if player is host/admin
	_ = player

	status := c.Query("status") // optional filter

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	suggestions, err := db.GetTileSuggestions(ctx, statusPtr)
	if err != nil {
		log.Printf("failed to get suggestions: %v", err)
		return utils.NewApiError("Failed to get suggestions", 0x0604).AsResponse(c)
	}

	return c.JSON(fiber.Map{"suggestions": suggestions})
}

// UpdateSuggestion handles PUT /api/suggestions/:id
func UpdateSuggestion(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	player, ok := c.Locals("player").(*models.Player)
	if !ok {
		return utils.NewApiError("Authentication required", 0x0401).AsResponse(c)
	}

	// TODO: Check if player is host/admin

	id := c.Params("id")
	if id == "" {
		return utils.NewApiError("Suggestion ID required", 0x0605).AsResponse(c)
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.NewApiError("Invalid request body", 0x0606).AsResponse(c)
	}

	if req.Status == "" {
		return utils.NewApiError("Status is required", 0x0607).AsResponse(c)
	}

	suggestion, err := db.UpdateTileSuggestion(ctx, id, req.Status, &player.ID)
	if err != nil {
		log.Printf("failed to update suggestion: %v", err)
		return utils.NewApiError("Failed to update suggestion", 0x0608).AsResponse(c)
	}

	if suggestion == nil {
		return utils.NewApiError("Suggestion not found", 0x0609).AsResponse(c)
	}

	return c.JSON(suggestion)
}
