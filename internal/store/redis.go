package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(url string) (*RedisStore, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid ALTCHA_REDIS_URL: %w", err)
	}
	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}
	return &RedisStore{client: client}, nil
}

func (s *RedisStore) Consume(nonce string, expiry time.Time) (bool, error) {
	ttl := time.Until(expiry)
	if ttl <= 0 {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// SET NX atomically sets the key only if it does not already exist.
	set, err := s.client.SetNX(ctx, nonce, 1, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis SetNX failed: %w", err)
	}
	return set, nil
}
