package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matoous/go-nanoid/v2"
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

	// Run migrations
	if err := RunMigrations(ctx); err != nil {
		log.Printf("database: failed to run migrations: %v", err)
		pool = nil
		p.Close()
		return
	}
}

// Pool returns the initialized pool or nil if not available.
func Pool() *pgxpool.Pool {
	return pool
}

// generateSessionID generates a 32-character session ID
func generateSessionID() string {
	id, _ := gonanoid.New(32)
	return id
}

// RunMigrations applies all pending migrations and seeds
func RunMigrations(ctx context.Context) error {
	if pool == nil {
		return fmt.Errorf("database not initialized")
	}

	// Get list of migration directories
	migrationsDir := "migrations"
	entries, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "00") {
			versions = append(versions, entry.Name())
		}
	}
	sort.Strings(versions)

	// Apply each migration
	for _, version := range versions {
		if err := applyMigration(ctx, filepath.Join(migrationsDir, version)); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", version, err)
		}
	}

	log.Println("All migrations applied successfully")
	return nil
}

func applyMigration(ctx context.Context, migrationDir string) error {
	version := filepath.Base(migrationDir)

	// Check if already applied
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}
	if count > 0 {
		log.Printf("Migration %s already applied, skipping", version)
		return nil
	}

	// Apply up.sql
	upPath := filepath.Join(migrationDir, "up.sql")
	upSQL, err := ioutil.ReadFile(upPath)
	if err != nil {
		return fmt.Errorf("failed to read up.sql: %w", err)
	}

	// Split by semicolon and execute each statement
	statements := strings.Split(string(upSQL), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	// Apply seed.sql
	seedPath := filepath.Join(migrationDir, "seed.sql")
	seedSQL, err := ioutil.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read seed.sql: %w", err)
	}

	statements = strings.Split(string(seedSQL), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute seed statement: %w", err)
		}
	}

	// Record migration as applied
	description := fmt.Sprintf("Applied migration %s", version)
	_, err = pool.Exec(ctx, "INSERT INTO schema_migrations (version, description) VALUES ($1, $2)", version, description)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	log.Printf("Applied migration %s", version)
	return nil
}

// RollbackMigration rolls back the last applied migration
func RollbackMigration(ctx context.Context) error {
	if pool == nil {
		return fmt.Errorf("database not initialized")
	}

	// Get the last applied migration
	var version string
	err := pool.QueryRow(ctx, "SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1").Scan(&version)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	migrationDir := filepath.Join("migrations", version)
	downPath := filepath.Join(migrationDir, "down.sql")
	downSQL, err := ioutil.ReadFile(downPath)
	if err != nil {
		return fmt.Errorf("failed to read down.sql: %w", err)
	}

	// Execute down.sql
	statements := strings.Split(string(downSQL), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute rollback statement: %w", err)
		}
	}

	// Remove from migrations table
	_, err = pool.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", version)
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	log.Printf("Rolled back migration %s", version)
	return nil
}
