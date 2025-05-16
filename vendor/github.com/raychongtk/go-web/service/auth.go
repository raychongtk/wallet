package service

import (
	"github.com/gin-gonic/gin"
	sessions "github.com/raychongtk/go-web/util"
)

func (s *Service) VerifySession() gin.HandlerFunc {
	return func(context *gin.Context) {
		sessionId := sessions.GetSessionId(context)

		if sessionId == nil {
			context.JSON(401, gin.H{
				"result": "unauthorized",
			})
			context.Abort()
		}
	}
}
