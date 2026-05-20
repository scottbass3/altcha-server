package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedStore struct {
	client *memcache.Client
}

func NewMemcachedStore(servers string) (*MemcachedStore, error) {
	addrs := strings.Split(servers, ",")
	for i, addr := range addrs {
		addrs[i] = strings.TrimSpace(addr)
	}
	client := memcache.New(addrs...)
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("memcached connection failed: %w", err)
	}
	return &MemcachedStore{client: client}, nil
}

func (s *MemcachedStore) Consume(nonce string, expiry time.Time) (bool, error) {
	ttl := int32(time.Until(expiry).Seconds())
	if ttl <= 0 {
		return false, nil
	}
	err := s.client.Add(&memcache.Item{
		Key:        nonce,
		Value:      []byte("1"),
		Expiration: ttl,
	})
	if err == memcache.ErrNotStored {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("memcached Add failed: %w", err)
	}
	return true, nil
}
