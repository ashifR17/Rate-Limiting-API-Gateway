package middleware

import (
	"net/http"

	"github.com/ashifR17/api-gateway/internal/service"
	"github.com/gin-gonic/gin"
)

func RateLimiter(limiter service.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ✅ Use request context (important)
		ctx := c.Request.Context()

		// Identify user
		userID := c.GetHeader("X-User-Id")
		if userID == "" {
			userID = c.ClientIP()
		}

		// Identify API
		apiID := c.FullPath()
		if apiID == "" {
			apiID = c.Request.URL.Path
		}

		// Build bucket keys
		globalKey := "rate:global"
		userKey := "rate:user:" + userID
		userAPIKey := "rate:user:" + userID + ":api:" + apiID

		// ✅ Pass context now
		allowed := limiter.Allow(
			ctx,
			globalKey,
			userKey,
			userAPIKey,
		)

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
