package service

import (
	"context"
	"math"
	"sync"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, global, user, userAPI string) bool
}

type tokenBucket struct {
	capacity     float64
	tokens       float64
	refillRate   float64
	lastrefillTs time.Time
}

type InMemoryLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
}

func NewInMemoryLimiter() *InMemoryLimiter {
	return &InMemoryLimiter{
		buckets: make(map[string]*tokenBucket),
	}
}

func (l *InMemoryLimiter) allowHelper(key string, capacity, refillRate float64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	bucket, exists := l.buckets[key]
	if !exists {
		bucket = &tokenBucket{
			capacity:     capacity,
			refillRate:   refillRate,
			tokens:       capacity,
			lastrefillTs: now,
		}
	}
	l.buckets[key] = bucket

	elapsed := now.Sub(bucket.lastrefillTs).Seconds()
	bucket.tokens = math.Min(
		capacity,
		bucket.tokens+(elapsed*bucket.refillRate),
	)

	bucket.lastrefillTs = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// Per request Allow criteria
// All 3 buckets should allow
func (l *InMemoryLimiter) Allow(global_key, userKey, userApiKey string) bool {

	//Global bucket limit
	if !l.allowHelper(global_key, 500, 50) {
		return false
	}

	//Per-User Bucket limit
	if !l.allowHelper(userKey, 100, 5) {
		return false
	}

	//Per-User-Per-API limit
	if !l.allowHelper(userApiKey, 20, 2) {
		return false
	}

	return true

}
