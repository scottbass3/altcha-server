package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/scottbass3/altcha-server/internal/config"
)

// Store tracks consumed challenge nonces to prevent replay attacks.
type Store interface {
	// Consume registers nonce and returns false if already used.
	Consume(nonce string, expiry time.Time) (bool, error)
}

func New(ctx context.Context, cfg config.Config) (Store, error) {
	switch strings.ToLower(cfg.StoreBackend) {
	case "redis":
		if cfg.RedisURL == "" {
			return nil, fmt.Errorf("ALTCHA_REDIS_URL is required when ALTCHA_STORE=redis")
		}
		return NewRedisStore(cfg.RedisURL)
	case "memcached":
		if cfg.MemcachedServers == "" {
			return nil, fmt.Errorf("ALTCHA_MEMCACHED_SERVERS is required when ALTCHA_STORE=memcached")
		}
		return NewMemcachedStore(cfg.MemcachedServers)
	case "memory", "":
		return NewMemoryStore(ctx), nil
	default:
		return nil, fmt.Errorf("unknown ALTCHA_STORE backend %q: must be memory, redis, or memcached", cfg.StoreBackend)
	}
}
