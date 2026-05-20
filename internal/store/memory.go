package store

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	mu      sync.Mutex
	entries map[string]time.Time
}

func NewMemoryStore(ctx context.Context) *MemoryStore {
	s := &MemoryStore{entries: make(map[string]time.Time)}
	go s.cleanup(ctx)
	return s
}

func (s *MemoryStore) Consume(nonce string, expiry time.Time) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.entries[nonce]; exists {
		return false, nil
	}
	s.entries[nonce] = expiry
	return true, nil
}

func (s *MemoryStore) cleanup(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			s.mu.Lock()
			for nonce, expiry := range s.entries {
				if now.After(expiry) {
					delete(s.entries, nonce)
				}
			}
			s.mu.Unlock()
		}
	}
}
