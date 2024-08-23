package cache

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"time"
)

type ICache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Del(key string) error
}

// Key wrapper key
func Key(key string) string {
	return fmt.Sprintf("%s:%s", utils.Getenv("APP_CODE", "gfly"), key)
}
