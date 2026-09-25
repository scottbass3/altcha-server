package store

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

const maxRelativeTTL = 30 * 24 * 60 * 60

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
	ttl := math.Ceil(time.Until(expiry).Seconds())
	if ttl <= 0 {
		return false, nil
	}
	exp := int32(ttl)
	// Memcached reads expirations above 30 days as absolute Unix timestamps.
	if ttl > maxRelativeTTL {
		exp = int32(expiry.Unix())
	}
	err := s.client.Add(&memcache.Item{
		Key:        keyPrefix + nonce,
		Value:      []byte("1"),
		Expiration: exp,
	})
	if err == memcache.ErrNotStored {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("memcached Add failed: %w", err)
	}
	return true, nil
}

func (s *MemcachedStore) Close() error { return s.client.Close() }
