package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
)

// Env var used for security between Go and Next.js
var nextAuthURL = os.Getenv("NEXT_AUTH_VALIDATE_URL") // e.g. http://localhost:3000/api/auth/validate
var internalSecret = os.Getenv("INTERNAL_API_SECRET") // optional shared secret

type validateResponse struct {
	Valid bool                   `json:"valid"`
	User  map[string]interface{} `json:"user"`
}

// getUserFromNext calls the Next.js /api/auth/validate endpoint.
func getUserFromNext(cookie string) (*validateResponse, error) {
	if nextAuthURL == "" {
		return nil, fmt.Errorf("NEXT_AUTH_VALIDATE_URL not configured")
	}

	req, err := http.NewRequest("GET", nextAuthURL, nil)
	if err != nil {
		return nil, err
	}

	// Pass the wos-session cookie along
	req.AddCookie(&http.Cookie{
		Name:  "wos-session",
		Value: cookie,
	})

	// Add shared secret header if configured
	if internalSecret != "" {
		req.Header.Set("X-Internal-Secret", internalSecret)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result validateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RequiredAuthMiddleware - require a valid WorkOS session
func RequiredAuthMiddleware(c *fiber.Ctx) error {
	sessionCookie := c.Cookies("wos-session")
	if sessionCookie == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing session cookie")
	}

	userInfo, err := getUserFromNext(sessionCookie)
	if err != nil || !userInfo.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired session")
	}

	c.Locals("user", userInfo.User)
	return c.Next()
}

// OptionalAuthMiddleware - optional session validation
func OptionalAuthMiddleware(c *fiber.Ctx) error {
	sessionCookie := c.Cookies("wos-session")
	if sessionCookie == "" {
		return c.Next()
	}

	userInfo, err := getUserFromNext(sessionCookie)
	if err == nil && userInfo.Valid {
		c.Locals("user", userInfo.User)
	}

	return c.Next()
}

// GetWorkOSUser retrieves the full user from WorkOS using the ID stored in ctx.Locals("user").
func GetWorkOSUser(c *fiber.Ctx) (*usermanagement.User, error) {
	raw := c.Locals("user")
	if raw != nil {
		if userMap, ok := raw.(map[string]interface{}); ok {
			if id, ok := userMap["id"].(string); ok && id != "" {
				user, err := usermanagement.DefaultClient.GetUser(c.Context(), usermanagement.GetUserOpts{
					User: id,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to fetch WorkOS user: %w", err)
				}
				return &user, nil
			}
			return nil, fmt.Errorf("missing or invalid user id in context")
		}
		return nil, fmt.Errorf("invalid user type: %T", raw)
	}
	return nil, fmt.Errorf("no user in context")
}
