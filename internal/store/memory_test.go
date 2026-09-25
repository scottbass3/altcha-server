package store

import (
	"testing"
	"time"
)

func TestMemoryStore(t *testing.T) {
	testStore(t, NewMemoryStore(t.Context()))
}

func TestMemoryStorePurge(t *testing.T) {
	s := NewMemoryStore(t.Context())
	now := time.Now()
	s.Consume("expired", now.Add(-time.Second))
	s.Consume("valid", now.Add(time.Minute))

	s.purge(now)

	if ok, _ := s.Consume("expired", now.Add(time.Minute)); !ok {
		t.Error("expired nonce was not purged")
	}
	if ok, _ := s.Consume("valid", now.Add(time.Minute)); ok {
		t.Error("valid nonce was purged")
	}
}
