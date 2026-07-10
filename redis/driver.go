package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gflydev/cache"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/redis/go-redis/v9"
)

// ========================================================================================
// 										Structure
// ========================================================================================

// Driver is a Redis-backed implementation of cache.ICache.
type Driver struct {
	redisCache *redis.Client
}

// Option customises the underlying redis.Options before the client is created.
type Option func(*redis.Options)

// WithPoolSize overrides the maximum number of socket connections.
func WithPoolSize(size int) Option {
	return func(o *redis.Options) { o.PoolSize = size }
}

// WithTimeouts overrides the dial, read and write timeouts.
func WithTimeouts(dial, read, write time.Duration) Option {
	return func(o *redis.Options) {
		o.DialTimeout = dial
		o.ReadTimeout = read
		o.WriteTimeout = write
	}
}

// New creates a Driver connected to the Redis server described by the
// REDIS_* environment variables. Extra Options may be supplied to tune the
// connection (pool size, timeouts, ...).
func New(opts ...Option) *Driver {
	// Define Redis database number.
	dbNumber := utils.Getenv("REDIS_DEFAULT_DB", 0)

	// Build Redis connection URL.
	redisConnURL := fmt.Sprintf(
		"%s:%d",
		utils.Getenv("REDIS_HOST", "localhost"),
		utils.Getenv("REDIS_PORT", 6379),
	)

	// Set Redis options.
	options := &redis.Options{
		Addr:     redisConnURL,
		Password: utils.Getenv("REDIS_PASSWORD", ""),
		DB:       dbNumber,
	}

	// Apply caller overrides.
	for _, opt := range opts {
		opt(options)
	}

	return &Driver{
		redisCache: redis.NewClient(options),
	}
}

// ========================================================================================
// 										Behaviour
// ========================================================================================

func (r *Driver) Set(key string, value interface{}, expiration time.Duration) error {
	if err := r.redisCache.Set(context.Background(), cache.Key(key), value, expiration).Err(); err != nil {
		log.Errorf("Error while writing Redis cache %q", err)
		return err
	}

	return nil
}

func (r *Driver) Get(key string) (interface{}, error) {
	val, err := r.redisCache.Get(context.Background(), cache.Key(key)).Result()
	if err != nil {
		// A missing key is an expected outcome, not a backend failure, so map it
		// to cache.ErrCacheMiss and keep it out of the error log.
		if errors.Is(err, redis.Nil) {
			return nil, cache.ErrCacheMiss
		}
		log.Errorf("Error while reading Redis key %q: %v", key, err)
		return nil, err
	}

	return val, nil
}

func (r *Driver) Del(key string) error {
	if err := r.redisCache.Del(context.Background(), cache.Key(key)).Err(); err != nil {
		log.Errorf("Error while deleting Redis key %q: %v", key, err)
		return err
	}

	return nil
}

// Ping verifies connectivity with the Redis server.
func (r *Driver) Ping() error {
	return r.redisCache.Ping(context.Background()).Err()
}

// Close releases the underlying Redis connection pool. It satisfies
// cache.ICloser so cache.Close() can shut the driver down cleanly.
func (r *Driver) Close() error {
	return r.redisCache.Close()
}
