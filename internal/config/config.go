package config

import (
	"fmt"
	"strconv"
	"time"
)

const (
	defaultHTTPAddr        = ":8080"
	defaultShutdownSeconds = 10
	minShutdownSeconds     = 1
	maxShutdownSeconds     = 300
)

type Config struct {
	HTTPAddr        string
	ShutdownTimeout time.Duration
}

// Load reads settings using getenv.
// We pass getenv in (instead of calling os.Getenv here) so tests
// can give fake values without touching the real environment.
func Load(getenv func(string) string) (Config, error) {
	addr := getenv("TASKBOARD_HTTP_ADDR")
	if addr == "" {
		addr = defaultHTTPAddr
	}

	seconds := defaultShutdownSeconds
	if raw := getenv("TASKBOARD_SHUTDOWN_TIMEOUT_SECONDS"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("TASKBOARD_SHUTDOWN_TIMEOUT_SECONDS must be a whole number: %w", err)
		}
		// Check the range BEFORE turning seconds into a Duration.
		// A huge number of times time.Second would overflow.
		if n < minShutdownSeconds || n > maxShutdownSeconds {
			return Config{}, fmt.Errorf("TASKBOARD_SHUTDOWN_TIMEOUT_SECONDS must be between %d and %d, got %d",
				minShutdownSeconds, maxShutdownSeconds, n)
		}
		seconds = n
	}

	return Config{
		HTTPAddr:        addr,
		ShutdownTimeout: time.Duration(seconds) * time.Second,
	}, nil
}
