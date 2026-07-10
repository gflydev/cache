package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/gflydev/core/utils"
)

// ========================================================================================
//                                     Errors
// ========================================================================================

// ErrNotRegistered is returned by the package-level helpers when they are used
// before a cache driver has been registered via Register.
var ErrNotRegistered = errors.New("cache: no driver registered, call cache.Register first")

// ErrCacheMiss is returned by Get when the requested key does not exist in the
// cache. Drivers should map their own "not found" sentinel (e.g. redis.Nil) to
// this error so callers can distinguish a genuine miss from a backend failure:
//
//	val, err := cache.Get(key)
//	if errors.Is(err, cache.ErrCacheMiss) {
//	    // key not present – recompute and Set it
//	}
var ErrCacheMiss = errors.New("cache: key not found")

// ========================================================================================
//                                     Struct
// ========================================================================================

// ICache is the contract every cache driver must satisfy.
type ICache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Del(key string) error
}

// ICloser is an optional interface a driver may implement to release resources
// (network connections, pools, ...). Close calls it when present.
type ICloser interface {
	Close() error
}

var cache ICache

// Register assigns the active cache manager. It is typically called once during
// application start-up, e.g. cache.Register(redis.New()).
func Register(c ICache) {
	cache = c
}

// Manager returns the currently registered cache driver, or nil if none has
// been registered yet.
func Manager() ICache {
	return cache
}

// ========================================================================================
//                                     Functions
// ========================================================================================

// Key builds a namespaced cache key by prefixing it with the APP_CODE value so
// that multiple applications can safely share a single backend.
func Key(key string) string {
	return fmt.Sprintf("%s:%s", utils.Getenv("APP_CODE", "gfly"), key)
}

// Set stores value under key for the given expiration. A zero expiration means
// the entry never expires.
func Set(key string, value interface{}, expiration time.Duration) error {
	if cache == nil {
		return ErrNotRegistered
	}
	return cache.Set(key, value, expiration)
}

// Get returns the value stored under key. It returns ErrCacheMiss when the key
// does not exist.
func Get(key string) (interface{}, error) {
	if cache == nil {
		return nil, ErrNotRegistered
	}
	return cache.Get(key)
}

// Del removes key from the cache. Deleting a missing key is not an error.
func Del(key string) error {
	if cache == nil {
		return ErrNotRegistered
	}
	return cache.Del(key)
}

// Close releases any resources held by the registered driver, if it implements
// ICloser. It is a no-op otherwise.
func Close() error {
	if cache == nil {
		return ErrNotRegistered
	}
	if closer, ok := cache.(ICloser); ok {
		return closer.Close()
	}
	return nil
}
