package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

const rateLimitConfigKey = "rate_limit:config"

type RateLimitConfig struct {
	GlobalCapacity  int `json:"global_capacity"`
	GlobalRate      int `json:"global_rate"` //refill rate
	UserCapacity    int `json:"user_capacity"`
	UserRate        int `json:"user_rate"`
	UserAPICapacity int `json:"user_api_capacity"`
	UserAPIRate     int `json:"user_api_rate"`
}

type ConfigService struct {
	client *redis.Client
}

func NewConfigService(client *redis.Client) *ConfigService {
	return &ConfigService{
		client: client,
	}
}

func (s *ConfigService) GetConfig(ctx context.Context) (*RateLimitConfig, error) {
	data, err := s.client.Get(ctx, rateLimitConfigKey).Result()

	if err == redis.Nil {
		//no config in Redis -> use default config
		return defaultConfig(), nil
	}

	if err != nil {
		return nil, err
	}

	var cfg RateLimitConfig
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (s *ConfigService) SetConfig(ctx context.Context, cfg *RateLimitConfig) error {
	if cfg == nil {
		return errors.New("set config can't be nil")
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, rateLimitConfigKey, data, 0).Err()
}

func defaultConfig() *RateLimitConfig {
	return &RateLimitConfig{
		GlobalCapacity:  500,
		GlobalRate:      50,
		UserCapacity:    100,
		UserRate:        5,
		UserAPICapacity: 20,
		UserAPIRate:     2,
	}
}
