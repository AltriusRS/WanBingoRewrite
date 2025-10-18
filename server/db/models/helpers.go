package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"time"
)

func (s *Show) HoursSince() float64 {
	now := time.Now().UTC()
	this := s.CreatedAt.UTC()

	diff := now.Sub(this)

	return diff.Hours()
}

// Hash generates a stable SHA-256 hash of the Show.
// Produces the same hash for identical struct content,
// regardless of map key order.
func (s *Show) Hash() (string, error) {
	// Copy to avoid modifying the original
	copy := *s

	// Normalize a Metadata map (sort keys for deterministic order)
	if copy.Metadata != nil {
		copy.Metadata = sortedMap(copy.Metadata)
	}

	data, err := json.Marshal(copy)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// Compare compares the current Show with another Show.
// Returns true if they are identical (based on hash), false if they differ.
func (s *Show) Compare(other *Show) (bool, error) {
	h1, err := s.Hash()
	if err != nil {
		return false, err
	}

	h2, err := other.Hash()
	if err != nil {
		return false, err
	}

	return h1 == h2, nil
}

// sortedMap returns a new map with deterministically sorted keys.
// Used to ensure stable JSON output for hashing.
func sortedMap(m map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make(map[string]interface{}, len(m))
	for _, k := range keys {
		sorted[k] = m[k]
	}
	return sorted
}
