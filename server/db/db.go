package db

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

// Init initializes a global pgx pool if DATABASE_URL is set.
// If the env var is missing or the connection fails, the function logs the error
// and leaves the pool as nil so the rest of the application can continue using
// file-based fallbacks.
func Init() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Println("database: DATABASE_URL not set; DB features disabled")
		return
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Printf("database: invalid DATABASE_URL: %v", err)
		return
	}
	// Reasonable defaults
	cfg.MaxConns = 4
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Printf("database: failed to create pool: %v", err)
		return
	}
	if err := p.Ping(ctx); err != nil {
		log.Printf("database: ping failed: %v", err)
		p.Close()
		return
	}
	pool = p
	log.Println("database: connected")
}

// Pool returns the initialized pool or nil if not available.
func Pool() *pgxpool.Pool {
	return pool
}

// generateID generates a random alphanumeric ID of specified length
func generateID(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// generateSessionID generates a 32-character session ID
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
