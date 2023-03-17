package builder

import (
	"os"

	"github.com/codecomet-io/go-core/log"
)

func init() {
	colors := Colors()
	err := os.Setenv("BUILDKIT_COLORS", colors)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set environment variable BUILDKIT_COLORS. Custom colors will not be used.")
	}
}