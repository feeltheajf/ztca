package log

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// supported colors
const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

// Config holds logging configuration
type Config struct {
	Level  string `yaml:"level"  bind:"required"`
	Format string `yaml:"format" bind:"required"`
}

// Setup initializes global logging subsystem
func Setup(cfg *Config) error {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("unknown log level: '%s'", cfg.Level)
	}
	zerolog.SetGlobalLevel(level)

	switch cfg.Format {
	case "json":
		break
	case "text":
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:         os.Stderr,
				FormatLevel: consoleFormatLevel,
			},
		)
	default:
		return fmt.Errorf("unknown log format: '%s'", cfg.Format)
	}
	return nil
}

// consoleFormatLevel is a custom log formatter for prettier text logs
func consoleFormatLevel(level interface{}) string {
	if levelString, ok := level.(string); ok {
		switch levelString {
		case "trace":
			return colorize("[T]", colorMagenta)
		case "debug":
			return colorize("[D]", colorCyan)
		case "info":
			return colorize("[I]", colorGreen)
		case "warn":
			return colorize("[W]", colorYellow)
		case "error":
			return colorize("[E]", colorRed)
		case "fatal":
			return colorize(colorize("[F]", colorRed), colorBold)
		case "panic":
			return colorize(colorize("[P]", colorRed), colorBold)
		}
	}
	return colorize("[?]", colorBold)
}

// colorize returns ANSI-colored strings
func colorize(message string, color int) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, message)
}
