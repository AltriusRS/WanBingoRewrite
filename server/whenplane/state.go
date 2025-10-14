package whenplane

import (
	"sync/atomic"
)

// aggregateCache stores the latest aggregate JSON payload as raw bytes.
var aggregateCache atomic.Value // holds []byte

func GetAggregateCache() Aggregate {
	return aggregateCache.Load().(Aggregate)
}

func UpdateAggregateCache(aggregate Aggregate) {
	aggregateCache.Store(aggregate)
}
