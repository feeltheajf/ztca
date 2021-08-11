package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds logging configuration
type Config struct {
	Level string `yaml:"level" bind:"required"`
}

// Setup initializes global logging subsystem
func Setup(cfg *Config) error {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out: os.Stderr,
		},
	)

	return nil
}
