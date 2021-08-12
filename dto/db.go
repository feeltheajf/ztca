package dto

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // GORM driver for sqlite3
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	levelDebug = "sql"
	levelInfo  = "log"
)

var (
	ctx zerolog.Logger

	all = []interface{}{
		&Certificate{},
	}

	db *gorm.DB
)

// Config holds database configuration
type Config struct {
	URL         string `yaml:"url" bind:"required"`
	Driver      string `yaml:"driver" bind:"required"`
	Automigrate bool   `yaml:"automigrate"`
}

// Setup initializes the database connection
func Setup(cfg *Config) (err error) {
	ctx = log.With().Str("module", "db").Logger()

	db, err = gorm.Open(cfg.Driver, cfg.URL)
	if err != nil {
		return err
	}

	db.LogMode(true)
	db.SetLogger(new(sqlLogger))
	db.DB().SetMaxOpenConns(1)

	if cfg.Automigrate {
		err := db.AutoMigrate(all...).Error
		if err != nil {
			return err
		}
	}

	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
	return nil
}

// sqlLogger implements gorm.Logger interface
type sqlLogger struct{}

func (l *sqlLogger) Print(values ...interface{}) {
	if len(values) < 3 {
		return
	}

	message := "sql"
	level := values[0].(string)
	path := values[1].(string)
	result := values[2]

	if level == levelDebug {
		ctx.Debug().
			Int64("elapsed_us", result.(time.Duration).Microseconds()).
			Str("sql", values[3].(string)).
			Interface("values", values[4]).
			Str("path", path).
			Msg(message)
		return
	}

	if level == levelInfo {
		switch result.(type) {
		case error:
			// errors are logged by middleware
			break
		default:
			ctx.Info().
				Interface("result", result).
				Msg(message)
		}
	}
}
