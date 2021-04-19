package api

import (
	"github.com/gin-gonic/gin"

	"github.com/feeltheajf/ztca/errdefs"
)

func handle(c *gin.Context, err error) {
	c.AbortWithStatusJSON(errdefs.GetStatusCode(err), err)
	c.Set(contextError, err)
}

func noRoute(c *gin.Context) {
	handle(c, errdefs.NotFound("URL not found"))
}
