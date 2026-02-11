package service

import (
	"context"
	_ "embed"
	"time"

	internalredis "github.com/ashifR17/api-gateway/internal/redis"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/token_bucket.lua
var luaScript string

type RedisLimiter struct {
	client        *redis.Client
	script        *redis.Script
	configService *ConfigService
}

func NewRedisLimiter(client *redis.Client, cfg *ConfigService) *RedisLimiter {
	return &RedisLimiter{
		client:        client,
		script:        redis.NewScript(luaScript),
		configService: cfg,
	}
}

// Atomic 3-bucket execution
func (l *RedisLimiter) Allow(ctx context.Context, global, user, userAPI string) bool {

	cfg, err := l.configService.GetConfig(ctx)
	if err != nil {
		return false
	}

	now := time.Now().Unix()

	res, err := l.script.Run(
		internalredis.Ctx,
		l.client,
		[]string{global, user, userAPI},

		// Global bucket
		cfg.GlobalCapacity, cfg.GlobalRate,
		// User bucket
		cfg.UserCapacity, cfg.UserRate,

		// User-API bucket
		cfg.UserAPICapacity, cfg.UserAPIRate,

		now, 3600,
	).Int()

	if err != nil {
		return false
	}

	return res == 1
}
