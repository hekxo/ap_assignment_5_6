package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func RateLimiter(duration time.Duration) gin.HandlerFunc {
	ticker := time.NewTicker(duration)
	return func(c *gin.Context) {
		select {
		case <-ticker.C:
			c.Next()
		default:
			log.Warn("Rate limit hit")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		}
	}
}
