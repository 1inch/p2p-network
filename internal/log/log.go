// Package log provides logging utilities.
package log

import "log/slog"

// ParseLevel parses log level from string value.
func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
