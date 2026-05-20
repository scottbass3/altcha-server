package store

import (
	"context"
	"sync"
	"time"
)

// NonceStore prevents replay attacks by tracking consumed challenge hashes.
type NonceStore struct {
	mu      sync.Mutex
	entries map[string]time.Time
}

func NewNonceStore(ctx context.Context) *NonceStore {
	ns := &NonceStore{entries: make(map[string]time.Time)}
	go ns.cleanup(ctx)
	return ns
}

// Consume registers the nonce and returns false if it was already used.
func (ns *NonceStore) Consume(nonce string, expiry time.Time) bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if _, exists := ns.entries[nonce]; exists {
		return false
	}
	ns.entries[nonce] = expiry
	return true
}

func (ns *NonceStore) cleanup(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			ns.mu.Lock()
			for nonce, expiry := range ns.entries {
				if now.After(expiry) {
					delete(ns.entries, nonce)
				}
			}
			ns.mu.Unlock()
		}
	}
}
