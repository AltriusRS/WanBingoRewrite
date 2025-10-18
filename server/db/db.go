package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"
	"wanshow-bingo/db/models"

	"github.com/jackc/pgx/v5"
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

func GetLatestShow(ctx context.Context) (*models.Show, error) {
	pool := Pool()

	latestShowRows, err := pool.Query(ctx, "SELECT * FROM shows ORDER BY scheduled_time DESC LIMIT 1")
	defer latestShowRows.Close()

	if err != nil {
		return nil, err
	}

	latestShows, err := pgx.CollectRows(latestShowRows, pgx.RowToStructByName[models.Show])

	if err != nil {
		return nil, err
	}

	latestShow := latestShows[0]

	return &latestShow, nil
}

func PersistShow(ctx context.Context, show *models.Show) error {
	if len(show.ID) < 10 {
		return errors.New("invalid show id")
	}

	pool := Pool()

	persisted, err := pool.Query(
		ctx,
		"UPDATE shows SET youtube_id = $1, scheduled_time = $2, actual_start_time = $3, thumbnail = $4, metadata = $5 WHERE id = $6",
		show.YoutubeID,
		show.ScheduledTime,
		show.ActualStartTime,
		show.Thumbnail,
		show.Metadata,
		show.ID,
	)

	defer persisted.Close()

	if err != nil {
		log.Printf("database: failed to update show: %v", err)
		return err
	}

	return nil
}

func SaveMessage(ctx context.Context, msg *models.Message) error {
	latestShow, err := GetLatestShow(ctx)
	if err != nil {
		log.Printf("[SSE ClientChannel] - Failed to retrieve latest show - %v", err)
		return err
	}

	pool := Pool()

	if pool == nil {
		log.Printf("[SSE ClientChannel] - Failed to connect to database")
		return err
	}

	msg.ShowID = latestShow.ID

	log.Printf("[SSE ClientChannel] - Saving message %v to DB", msg)

	res, err := pool.Query(
		ctx,
		"INSERT INTO messages (id, show_id, player_id, contents, system, replying) VALUES ($1, $2, $3, $4, $5, $6)",
		msg.ID, msg.ShowID, msg.PlayerID, msg.Contents, msg.System, msg.Replying,
	)

	defer res.Close()

	if err != nil {
		log.Printf("database: failed to save message: %v", err)
		return err
	}

	if res.Err() != nil {
		log.Printf("database: failed to save message: %v", res.Err())
	}

	return nil
}

func GetMessageHistory(ctx context.Context) ([]models.Message, error) {
	latestShow, err := GetLatestShow(ctx)
	if err != nil {
		return nil, err
	}

	pool := Pool()

	res, err := pool.Query(ctx, "SELECT * FROM messages WHERE show_id = $1 ORDER BY created_at DESC LIMIT 30", latestShow.ID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	// RowToStructByName expects the struct type, not pointer
	messages, err := pgx.CollectRows(res, pgx.RowToStructByName[models.Message])
	if err != nil {
		return nil, err
	}

	return messages, nil
}
