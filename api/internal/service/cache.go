package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
}

func NewCacheService(addr string) *CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &CacheService{client: rdb}
}

func (s *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(ctx, key, value, expiration).Err()
}

func (s *CacheService) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}
