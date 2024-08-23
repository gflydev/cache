package cache

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"time"
)

// ========================================================================================
//                                     Struct
// ========================================================================================

type ICache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Del(key string) error
}

var cache ICache

// Register assign cache manager
func Register(c ICache) {
	cache = c
}

// ========================================================================================
//                                     Functions
// ========================================================================================

// Key wrapper key
func Key(key string) string {
	return fmt.Sprintf("%s:%s", utils.Getenv("APP_CODE", "gfly"), key)
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return cache.Set(key, value, expiration)
}

func Get(key string) (interface{}, error) {
	return cache.Get(key)
}

func Del(key string) error {
	return cache.Del(key)
}
