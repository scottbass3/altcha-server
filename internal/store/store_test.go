package store

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func testStore(t *testing.T, s Store) {
	t.Helper()
	prefix := fmt.Sprintf("test-%d-", time.Now().UnixNano())

	t.Run("consume once", func(t *testing.T) {
		expiry := time.Now().Add(time.Minute)
		if ok, err := s.Consume(prefix+"once", expiry); err != nil || !ok {
			t.Fatalf("first Consume = %v, %v; want true", ok, err)
		}
		if ok, err := s.Consume(prefix+"once", expiry); err != nil || ok {
			t.Fatalf("second Consume = %v, %v; want false", ok, err)
		}
	})

	t.Run("concurrent consume", func(t *testing.T) {
		var wins atomic.Int32
		var wg sync.WaitGroup
		for range 50 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if ok, err := s.Consume(prefix+"concurrent", time.Now().Add(time.Minute)); err == nil && ok {
					wins.Add(1)
				}
			}()
		}
		wg.Wait()
		if wins.Load() != 1 {
			t.Fatalf("%d concurrent Consume calls succeeded; want 1", wins.Load())
		}
	})

	t.Run("long ttl", func(t *testing.T) {
		expiry := time.Now().Add(60 * 24 * time.Hour)
		if ok, err := s.Consume(prefix+"long", expiry); err != nil || !ok {
			t.Fatalf("first Consume = %v, %v; want true", ok, err)
		}
		if ok, err := s.Consume(prefix+"long", expiry); err != nil || ok {
			t.Fatalf("second Consume = %v, %v; want false", ok, err)
		}
	})
}

func TestRedisStore(t *testing.T) {
	url := os.Getenv("ALTCHA_TEST_REDIS_URL")
	if url == "" {
		t.Skip("ALTCHA_TEST_REDIS_URL not set")
	}
	s, err := NewRedisStore(url)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	testStore(t, s)
}

func TestMemcachedStore(t *testing.T) {
	servers := os.Getenv("ALTCHA_TEST_MEMCACHED")
	if servers == "" {
		t.Skip("ALTCHA_TEST_MEMCACHED not set")
	}
	s, err := NewMemcachedStore(servers)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	testStore(t, s)
}
