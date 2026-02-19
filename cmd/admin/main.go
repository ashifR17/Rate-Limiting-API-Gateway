package main

import (
	"net/http"

	internalredis "github.com/ashifR17/api-gateway/internal/redis"
	"github.com/ashifR17/api-gateway/internal/service"
	"github.com/gin-gonic/gin"
)

// Dummy admin token
const ADMIN_TOKEN = "supersecret"

// Middleware: dummy admin auth
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Admin-Key") != ADMIN_TOKEN {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func main() {
	rdb := internalredis.New()
	ctx := internalredis.Ctx
	configService := service.NewConfigService(rdb)

	r := gin.Default()
	r.Static("/static", "./static")
	r.StaticFile("/admin", "./admin.html")

	admin := r.Group("/admin", AdminAuth())

	admin.POST("/config", func(c *gin.Context) {
		var cfg service.RateLimitConfig
		if err := c.ShouldBindJSON(&cfg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body"})
			return
		}
		if err := configService.SetConfig(ctx, &cfg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save config"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Rate limit config updated"})

	})

	admin.GET("/config", func(c *gin.Context) {
		cfg, err := configService.GetConfig(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch config"})

			return
		}
		c.JSON(http.StatusOK, cfg)
	})

	r.Run(":9090")

}
