package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/feeltheajf/ztca/dto"
)

type session struct {
	UUID string
	User *dto.User
}

func getSession(c *gin.Context) *session {
	return &session{
		UUID: uuid.NewString(),
		User: &dto.User{
			Name: "test client",
		},
	}
}
