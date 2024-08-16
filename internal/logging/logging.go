package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger is the global logger.
var Logger zerolog.Logger

// SetupLogger configures the global logger.
func SetupLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"})
	log.Logger = Logger
}

// ColoredURL returns the URL wrapped in ANSI escape codes for coloring.
func ColoredURL(url string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", url) // Light blue color
}
