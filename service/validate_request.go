package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *Service) ValidateRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing X-Request-ID header",
			})
			c.Abort()
			return
		}

		exists := s.memoryStore.Exists(context.Background(), requestID).Val()
		if exists > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Duplicate request",
			})
			c.Abort()
			return
		}
		s.memoryStore.Set(context.Background(), requestID, "true", time.Hour)
		c.Next()
	}
}
