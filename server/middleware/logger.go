package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestLogger is a Fiber middleware that logs incoming HTTP requests
// together with the time taken to serve them.
func RequestLogger(c *fiber.Ctx) error {
	// 1 Capture start time
	start := time.Now()

	// 2 Let the next handlers run
	err := c.Next() // <--- this is where the request actually gets processed

	// 3 Measure elapsed time
	duration := time.Since(start)

	// 4 Gather some useful bits of data
	method := string(c.Method())            // GET, POST, …
	path := c.Path()                        // "/chat/"
	statusCode := c.Response().StatusCode() // 200, 404, …

	// Optional: include query string (useful for debugging)
	if q := c.Request().URI().QueryArgs(); q.Len() > 0 {
		path += "?" + q.String()
	}

	// 5 Log it.  Replace fmt.Printf with whatever logger you like.
	fmt.Printf("%s %s %s → %d (%v)\n", time.Now().Format("2006-01-02 15:04:05"), method, path, statusCode, duration)

	return err
}
