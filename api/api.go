package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	// TODO use values from ztunnel/x/http
	headerAPIToken  = "X-Api-Token" // #nosec: G101
	headerRequestID = "X-Request-Id"
	headerUserAgent = "User-Agent"

	contextError = "error"
)

var (
	server *http.Server
	config *Config
)

type Config struct {
	HTTP *HTTP `yaml:"http"`
	Auth *Auth `yaml:"auth"`
}

type HTTP struct {
	Address string `yaml:"address"`
}

type Auth struct {
	APIToken string `yaml:"apiToken"`
}

func Serve() error {
	// TODO TLS support
	log.Info().
		Str("address", server.Addr).
		Msg("running HTTP server")
	return server.ListenAndServe()
}

func Setup(cfg *Config) error {
	server = &http.Server{
		Addr:         cfg.HTTP.Address,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
	}
	config = cfg

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
	v1.Use(hasToken)
	{
		v1.POST("/requests/yubikey", yubikey)

		v1.GET("/requests/nonce", nonce)
		v1.POST("/requests/generate", generate)
		v1.POST("/requests/activate", activate)
	}

	server.Handler = r
	return nil
}
