package watcher

import (
	"context"
	"log"
	"strings"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"
	"wanshow-bingo/sse"
	"wanshow-bingo/utils"
	"wanshow-bingo/whenplane"

	"github.com/jackc/pgx/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

var AggregateChan chan *whenplane.Aggregate

func init() {
	AggregateChan = make(chan *whenplane.Aggregate, 100)
	AggregateHandler()
}

func AggregateHandler() {
	go func() {
		log.Println("Starting aggregate handler")
		for event := range AggregateChan {
			HandleAggregateEvent(event)
			time.Sleep(time.Second)
		}
	}()
}

func HandleAggregateEvent(newAggregate *whenplane.Aggregate) {
	oldShow, err := whenplane.GetAggregateCache()

	if err != nil {
		log.Println("Error getting current aggregate cache")
		log.Println("Creating new aggregate")
		show := BuildShowFromAggregate(newAggregate)
		whenplane.SetAggregateCache(show)
		log.Println("New aggregate created")
		return
	}

	newShow := BuildShowFromAggregate(newAggregate)
	if newShow == nil {
		log.Printf("Failed to build a show from the aggregate provided (new) - %v", newShow)
		return
	}

	// Update show state based on aggregate
	err = UpdateShowState(context.Background(), newAggregate)
	if err != nil {
		log.Printf("Failed to update show state: %v", err)
	}

	showsAreEqual, err := oldShow.Compare(newShow)

	if err != nil {
		log.Printf("Failed to compare shows diff from the old show - %v", oldShow)
	}

	if !showsAreEqual {
		log.Printf("[AGGREGATE] The shows are different")
		existingShow, err := db.GetLatestShow(context.Background())

		if err != nil {
			log.Printf("[AGGREGATE] Failed to get latest show - %v", err)
		}

		whenplane.SetAggregateCache(newShow)

		if existingShow.Metadata["title"] != newShow.Metadata["title"] {
			log.Printf("[AGGREGATE] The database show is different from the new show")

			ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
			defer cancel()

			showTitle := ExtractShowTitle(existingShow.Metadata["title"].(string))
			distanceFromHello := Distances(showTitle, "Hello, Floatplane!")
			distanceFromTitle := Distances(showTitle, existingShow.Metadata["title"].(string))

			if distanceFromHello < 5 && distanceFromTitle > 5 {
				// The title is probably set to "Hello, Floaptlane!" by Dan arriving on set
				// We should use this as the marker to lock the tiles in for the current show
				// (and probably send a system message)

				TitleChanged(ctx, newShow)
			} else {
				// The title has been set to (probably) it's final form
				newShow.ID = existingShow.ID

				if newAggregate.Youtube.IsLive && newShow.ActualStartTime == nil {
					startTimeApprox := time.Now().UTC()
					newShow.ActualStartTime = &startTimeApprox
				}

				err := db.PersistShow(ctx, newShow)

				if err != nil {
					log.Printf("[AGGREGATE] Failed to persist show - %v", err)
				}

				//fpMeta := newShow.Metadata["floatplane"].(map[string]interface{})

				//if fpMeta[""]
			}

		}
	}

	err = BroadcastToHubs(newShow)

	if err != nil {
		log.Printf("[AGGREGATE] Error broadcasting to hubs: %v", err)
	}
}

func BroadcastToHubs(payload *models.Show) error {
	utils.Debugln("[AGGREGATE] Broadcasting aggregate state to hubs")

	chatHub := sse.GetChatHub()
	if chatHub != nil {
		utils.Debugln("[AGGREGATE] Broadcasting aggregate state to chat hub")
		chatHub.BroadcastEvent("whenplane.aggregate", payload)
	}

	hostHub := sse.GetHostHub()
	if hostHub != nil {
		utils.Debugln("[AGGREGATE] Broadcasting aggregate state to host hub")
		hostHub.BroadcastEvent("whenplane.aggregate", payload)
	}

	return nil
}

func BuildShowFromAggregate(aggregate *whenplane.Aggregate) *models.Show {
	next, err := NextWAN()

	if err != nil {
		log.Printf("[AGGREGATE] Failed to build show from aggregate: %v", err)
		return nil
	}

	var thumbnail *string

	if aggregate.Youtube.VideoID != nil {
		temp := "https://i.ytimg.com/vi/" + *aggregate.Youtube.VideoID + "/maxresdefault.jpg"
		thumbnail = &temp
	} else {
		thumbnail = &aggregate.Floatplane.Thumbnail
	}

	metadata := make(map[string]interface{})

	metadata["title"] = ExtractShowTitle(aggregate.Floatplane.Title)

	var actualStartTime *time.Time

	if aggregate.Youtube.IsLive {
		start := time.Now().UTC()
		actualStartTime = &start
	}

	floatplaneMetadata := map[string]interface{}{
		"thumbnail":        thumbnail,
		"is_live":          aggregate.Floatplane.IsLive,
		"is_wan":           aggregate.Floatplane.IsWAN,
		"is_thumbnail_new": aggregate.Floatplane.IsThumbnailNew,
		"title":            aggregate.Floatplane.Title,
	}

	youtubeMetadata := map[string]interface{}{
		"title":    aggregate.Floatplane.Title,
		"video_id": aggregate.Youtube.VideoID,
		"upcoming": aggregate.Youtube.Upcoming,
		"is_live":  aggregate.Floatplane.IsLive,
	}

	metadata["floatplane"] = floatplaneMetadata
	metadata["youtube"] = youtubeMetadata

	show := &models.Show{
		ID:              "",
		State:           models.ShowStateScheduled,
		YoutubeID:       aggregate.Youtube.VideoID,
		ScheduledTime:   &next,
		ActualStartTime: actualStartTime,
		Thumbnail:       thumbnail,
		Metadata:        metadata,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		DeletedAt:       nil,
	}

	return show
}

func TitleChanged(ctx context.Context, newShow *models.Show) {
	log.Println("[AGGREGATE] Title changed")

	latestShow, err := db.GetLatestShow(ctx)

	if err != nil {
		log.Printf("[AGGREGATE] DB: Failed to get latest show - %v", err)
		return
	}

	log.Printf("[AGGREGATE] Latest show found - %s", latestShow.ID)

	diff := latestShow.HoursSince()

	if diff > 110 {
		pool := db.Pool()
		log.Printf("[AGGREGATE] New show is different to current show")
		// Start a new transaction to minimus risk if the query fails
		tx, err := pool.Begin(ctx)

		if err != nil {
			log.Printf("[AGGREGATE] Error creating database transaction: %v", err)
			return
		} else {
			err := NewShowProtocol(ctx, tx, newShow)
			if err != nil {
				log.Printf("[AGGREGATE] Error running show protocol: %v", err)
				err = tx.Rollback(ctx)
				if err != nil {
					log.Printf("[AGGREGATE] Error rolling back database transaction: %v", err)
				}
				return
			}
		}

	} else {
		log.Printf("[AGGREGATE] New show is generated too close to existing show - %f / 110 hours differences", diff)
	}
}

func NewShowProtocol(ctx context.Context, tx pgx.Tx, newShow *models.Show) error {
	// Generate a new show entry in the transaction
	showId, err := CreateNewShow(ctx, tx, newShow)

	if err != nil {
		log.Printf("[AGGREGATE] Error creating show: %v", err)
		return err
	}

	// Insert was successful. Now we can make the tiles
	playingField, err := GenerateRandomPlayingField(ctx, tx, showId)

	if err != nil {
		log.Printf("[AGGREGATE] Error generating playing field: %v", err)
		return err
	}

	log.Printf("[AGGREGATE] Playing field: %s", playingField)

	err = tx.Commit(ctx)

	if err != nil {
		log.Printf("[AGGREGATE] Error committing database transaction: %v", err)
		return err
	}

	return nil
}

func NextWAN() (time.Time, error) {
	loc, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now().In(loc)

	// Calculate days until next Friday (Weekday == 5)
	daysUntilFriday := (int(time.Friday) - int(now.Weekday()) + 7) % 7
	if daysUntilFriday == 0 {
		// If today is Friday but it's past 4:30 PM, go to next week
		if now.Hour() > 16 || (now.Hour() == 16 && now.Minute() >= 30) {
			daysUntilFriday = 7
		}
	}

	nextFriday := now.AddDate(0, 0, daysUntilFriday)
	target := time.Date(
		nextFriday.Year(),
		nextFriday.Month(),
		nextFriday.Day(),
		16, 30, 0, 0,
		loc,
	)

	return target, nil
}

func CreateNewShow(ctx context.Context, tx pgx.Tx, newShow *models.Show) (*string, error) {
	id, err := gonanoid.New(10)

	if err != nil {
		log.Printf("[AGGREGATE] NID: Failed to generate id - %s", err)
		return nil, err
	}

	newShow.ID = id

	insertResult, err := tx.Query(
		ctx,
		"INSERT INTO shows (id, state, youtube_id, scheduled_time, actual_start_time, thumbnail, metadata, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		newShow.ID,
		string(newShow.State),
		newShow.YoutubeID,
		newShow.ScheduledTime,
		newShow.ActualStartTime,
		newShow.Thumbnail,
		newShow.Metadata,
		newShow.CreatedAt,
		newShow.UpdatedAt,
		newShow.DeletedAt,
	)

	if err != nil {
		log.Printf("[AGGREGATE] DB: Failed to insert show - %s", err)

		return nil, err
	}

	insertResult.Close()

	return &newShow.ID, nil
}

func GenerateRandomPlayingField(ctx context.Context, tx pgx.Tx, showId *string) (*[]string, error) {
	randomTilesQuery, err := tx.Query(ctx, "SELECT * FROM tiles ORDER BY random() LIMIT 90")

	if err != nil {
		log.Printf("[AGGREGATE] DB: Failed to generate random playing field - %s", err)
		return nil, err
	}

	tiles, err := pgx.CollectRows(randomTilesQuery, pgx.RowToStructByName[models.Tile])

	if err != nil {
		log.Printf("[AGGREGATE] DB: Failed to deserialize random playing field - %s", err)
		return nil, err
	}

	randomTilesQuery.Close()

	showTiles := make([]models.ShowTile, 0, len(tiles))
	tileIDs := make([]string, 0, len(tiles))

	for _, tile := range tiles {
		tileIDs = append(tileIDs, tile.ID)
		showTiles = append(showTiles, models.ShowTile{
			ShowID:    *showId,
			TileID:    tile.ID,
			Weight:    tile.Weight,
			Score:     tile.Score,
			CreatedAt: tile.CreatedAt,
			UpdatedAt: tile.UpdatedAt,
			DeletedAt: nil,
		})
	}

	rows := make([][]any, len(showTiles))
	for i, t := range showTiles {
		rows[i] = []any{
			t.ShowID,
			t.TileID,
			t.Weight,
			t.Score,
			t.CreatedAt,
			t.UpdatedAt,
			t.DeletedAt,
		}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"show_tiles"},
		[]string{"show_id", "tile_id", "weight", "score", "created_at", "updated_at", "deleted_at"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		log.Printf("[AGGREGATE] DB: Failed to bulk insert show tiles - %s", err)
		return nil, err
	}

	// Return all tile IDs for later board generation
	return &tileIDs, nil
}

func ExtractShowTitle(title string) string {
	parts := strings.SplitN(title, " - ", 2)

	if len(parts) < 2 {
		return title
	}
	return parts[0]
}

func Distances(a string, b string) int {
	return levenshtein.DistanceForStrings([]rune(a), []rune(b), levenshtein.DefaultOptions)
}

func UpdateShowState(ctx context.Context, aggregate *whenplane.Aggregate) error {
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		return err
	}

	newState := latestShow.State

	if aggregate.Youtube.IsLive {
		newState = models.ShowStateLive
	} else if aggregate.Youtube.VideoID != nil || aggregate.Floatplane.IsThumbnailNew {
		newState = models.ShowStateUpcoming
	} else if latestShow.State == models.ShowStateLive {
		newState = models.ShowStateFinished
	}

	if newState != latestShow.State {
		log.Printf("[AGGREGATE] Updating show state from %s to %s", latestShow.State, newState)
		err = db.UpdateShowState(ctx, latestShow.ID, newState)
		if err != nil {
			return err
		}

		// Handle WAN timer based on state change
		if newState == models.ShowStateLive && latestShow.State != models.ShowStateLive {
			// Show went live, create 4-hour timer
			timerID, _ := gonanoid.New(10)
			now := time.Now()
			expiresAt := now.Add(4 * time.Hour)
			timer := &models.Timer{
				ID:        timerID,
				Title:     "WAN Show Timer",
				Duration:  14400, // 4 hours in seconds
				CreatedBy: nil,   // System timer
				ShowID:    &latestShow.ID,
				StartsAt:  &now,
				ExpiresAt: &expiresAt,
				IsActive:  true,
				Settings:  map[string]interface{}{},
				CreatedAt: now,
				UpdatedAt: now,
			}
			err = db.PersistTimer(context.Background(), timer)
			if err != nil {
				log.Printf("[AGGREGATE] Failed to create WAN timer: %v", err)
			} else {
				log.Printf("[AGGREGATE] Created WAN timer for show %s", latestShow.ID)
			}
		} else if newState != models.ShowStateLive && latestShow.State == models.ShowStateLive {
			// Show went offline, cancel WAN timer
			err = db.StopActiveTimersByTitle(context.Background(), "WAN Show Timer", latestShow.ID)
			if err != nil {
				log.Printf("[AGGREGATE] Failed to stop WAN timer: %v", err)
			} else {
				log.Printf("[AGGREGATE] Stopped WAN timer for show %s", latestShow.ID)
			}
		}
	}

	return nil
}
