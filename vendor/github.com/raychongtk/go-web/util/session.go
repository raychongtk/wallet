package util

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateSession(ctx *gin.Context, username string) {
	session := sessions.Default(ctx)
	session.Set("SessionId", uuid.New().String())
	session.Set("Username", username)
	session.Save()
}

func GetSessionId(ctx *gin.Context) interface{} {
	session := sessions.Default(ctx)
	sessionId := session.Get("SessionId")
	return sessionId
}

func GetUsername(ctx *gin.Context) string {
	session := sessions.Default(ctx)
	userId := fmt.Sprintf("%v", session.Get("Username"))
	return userId
}
