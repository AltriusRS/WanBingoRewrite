package utils

import (
	"log"
	"os"
	"strings"
)

var debugEnabled bool

func init() {
	// Respect either WANBINGO_DEBUG or DEBUG env vars (case-insensitive, typical truthy values)
	debugEnabled = parseBoolEnv("WANBINGO_DEBUG") || parseBoolEnv("DEBUG")
}

func parseBoolEnv(key string) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch v {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func DebugEnabled() bool { return debugEnabled }

func Debugf(format string, args ...any) {
	if debugEnabled {
		log.Printf(format, args...)
	}
}

func Debugln(v ...any) {
	if debugEnabled {
		log.Println(v...)
	}
}
