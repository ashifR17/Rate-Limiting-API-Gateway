package main

import (
	"log"
	"os"

	"github.com/ashifR17/api-gateway/internal/handler"
	"github.com/ashifR17/api-gateway/internal/middleware"
	"github.com/ashifR17/api-gateway/internal/redis"
	"github.com/ashifR17/api-gateway/internal/service"

	"github.com/gin-gonic/gin"
)

/*HTTP Request
     ↓
Gin Middleware
     ↓
Limiter.Allow(ctx, ...)
     ↓
ConfigService
     ↓
Redis
     ↓
Atomic Lua Script
*/

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "gateway up",
		})
	})

	// limiter := service.NewInMemoryLimiter()
	redisClient := redis.New()
	configService := service.NewConfigService(redisClient)

	limiter := service.NewRedisLimiter(redisClient, configService)

	// Rate limit
	r.Use(middleware.RateLimiter(limiter))

	r.Any("api/*proxyPath", handler.ReverseProxy("http://localhost:9000"))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	log.Println("API Gateway runnning on :%v", port)

	r.Run(":" + port)
}
