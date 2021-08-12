package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	headerRequestID = "X-Request-Id"
	headerUserAgent = "User-Agent"

	contextError = "error"
)

var (
	ctx    zerolog.Logger
	server *http.Server
)

type Config struct {
	HTTP *HTTP `yaml:"http"`
}

type HTTP struct {
	Address string `yaml:"address"`
}

func Serve() error {
	ctx.Info().
		Str("address", server.Addr).
		Msg("running HTTP server")
	return server.ListenAndServe()
}

func Setup(cfg *Config) error {
	ctx = log.With().Str("module", "api").Logger()

	server = &http.Server{
		Addr:         cfg.HTTP.Address,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
	}

	if err := setupRoutes(cfg.HTTP); err != nil {
		return err
	}

	return nil
}

func setupRoutes(cfg *HTTP) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.NoRoute(noRoute)
	r.Use(recovery)
	r.Use(logging)

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		v1.POST("/certs/:username", issueCertificate)
		v1.DELETE("/certs/:username", revokeCertificate)
	}

	server.Handler = r
	return nil
}
