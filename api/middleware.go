package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/feeltheajf/ztca/errdefs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(error)
			if !ok {
				e = fmt.Errorf("%+v", err)
			}
			log.Error().
				Err(e).
				Bytes("stack", debug.Stack()).
				Msg("panic recovered")
		}
	}()
	c.Next()
}

func logging(c *gin.Context) {
	now := time.Now()
	c.Next()
	code := c.Writer.Status()
	ctx := log.With().
		Int("code", code).
		Str("method", c.Request.Method).
		Str("path", c.Request.RequestURI).
		Str("request_id", uuid.NewString()).
		Int64("elapsed_us", time.Since(now).Microseconds()).
		Str("user_ip", c.ClientIP()).
		Str("user_agent", c.GetHeader(headerUserAgent))

	s := getSession(c)
	if s != nil {
		ctx = ctx.
			Str("session_id", s.UUID).
			Str("user", s.User.Name)
	}

	if v, ok := c.Get(contextError); ok {
		err := v.(error)
		ctx = ctx.Err(err)
	}

	logger := ctx.Logger()
	var event *zerolog.Event
	switch {
	case code < http.StatusBadRequest:
		event = logger.Info()
	case code < http.StatusInternalServerError:
		event = logger.Warn()
	default:
		event = logger.Error()
	}
	event.Msg("request")
}

func hasToken(c *gin.Context) {
	apiToken := c.GetHeader(headerAPIToken)
	if apiToken == "" || apiToken != config.Auth.APIToken {
		handle(c, errdefs.Forbidden("invalid API token"))
	}
}
