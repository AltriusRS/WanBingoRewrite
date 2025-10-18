package whenplane

import (
	"errors"
	"sync"
	"wanshow-bingo/db/models"
)

var (
	mu             sync.RWMutex
	aggregateCache *models.Show
)

func SetAggregateCache(a *models.Show) {
	mu.Lock()
	defer mu.Unlock()
	aggregateCache = a
}

func GetAggregateCache() (*models.Show, error) {
	mu.RLock()
	defer mu.RUnlock()
	if aggregateCache == nil {
		return nil, errors.New("aggregate cache empty")
	}
	return aggregateCache, nil
}
