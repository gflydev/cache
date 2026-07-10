package redis

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestWithPoolSize(t *testing.T) {
	opts := &redis.Options{}
	WithPoolSize(42)(opts)
	if opts.PoolSize != 42 {
		t.Fatalf("PoolSize: want 42, got %d", opts.PoolSize)
	}
}

func TestWithTimeouts(t *testing.T) {
	opts := &redis.Options{}
	WithTimeouts(time.Second, 2*time.Second, 3*time.Second)(opts)

	if opts.DialTimeout != time.Second {
		t.Fatalf("DialTimeout: want 1s, got %v", opts.DialTimeout)
	}
	if opts.ReadTimeout != 2*time.Second {
		t.Fatalf("ReadTimeout: want 2s, got %v", opts.ReadTimeout)
	}
	if opts.WriteTimeout != 3*time.Second {
		t.Fatalf("WriteTimeout: want 3s, got %v", opts.WriteTimeout)
	}
}

func TestNewAppliesOptions(t *testing.T) {
	d := New(WithPoolSize(7))
	if d == nil || d.redisCache == nil {
		t.Fatal("New returned an incomplete Driver")
	}
	if got := d.redisCache.Options().PoolSize; got != 7 {
		t.Fatalf("PoolSize override not applied: got %d", got)
	}
	// The client is created lazily, so closing it should never error.
	if err := d.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
}
