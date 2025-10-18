package db

import (
	"context"
	"fmt"
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
	entries, err := os.ReadDir(migrationsDir)
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
		// If the schema_migrations table doesn't exist, assume this is the first run
		// and proceed with applying migrations (the table will be created by 001_core)
		if strings.Contains(err.Error(), "does not exist") {
			log.Printf("schema_migrations table does not exist, assuming first run")
		} else {
			return fmt.Errorf("failed to check migration status: %w", err)
		}
	} else if count > 0 {
		log.Printf("Migration %s already applied, skipping", version)
		return nil
	}

	// Begin transaction for atomic migration
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Apply up.sql
	upPath := filepath.Join(migrationDir, "up.sql")
	upSQL, err := os.ReadFile(upPath)
	if err != nil {
		return fmt.Errorf("failed to read up.sql: %w", err)
	}

	// Execute the entire up.sql file in a single transaction
	if _, err = tx.Exec(ctx, string(upSQL)); err != nil {
		return fmt.Errorf("failed to execute up.sql: %w", err)
	}

	// Execute the entire up.sql file in a single transaction
	if _, err = tx.Exec(ctx, string(upSQL)); err != nil {
		return fmt.Errorf("failed to execute up.sql: %w", err)
	}

	// Apply seed.sql
	seedPath := filepath.Join(migrationDir, "seed.sql")
	seedSQL, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read seed.sql: %w", err)
	}

	// Execute the entire seed.sql file in a single transaction
	if _, err = tx.Exec(ctx, string(seedSQL)); err != nil {
		return fmt.Errorf("failed to execute seed.sql: %w", err)
	}

	// Record migration as applied
	description := fmt.Sprintf("Applied migration %s", version)
	_, err = tx.Exec(ctx, "INSERT INTO schema_migrations (version, description) VALUES ($1, $2)", version, description)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
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

	// Begin transaction for atomic rollback
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin rollback transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	migrationDir := filepath.Join("migrations", version)
	downPath := filepath.Join(migrationDir, "down.sql")
	downSQL, err := os.ReadFile(downPath)
	if err != nil {
		return fmt.Errorf("failed to read down.sql: %w", err)
	}

	// Execute down.sql in transaction
	if _, err = tx.Exec(ctx, string(downSQL)); err != nil {
		return fmt.Errorf("failed to execute down.sql: %w", err)
	}

	// Remove from migrations table
	_, err = tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", version)
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	log.Printf("Rolled back migration %s", version)
	return nil
}
