package redis

import (
	"context"
	"fmt"
	"github.com/gflydev/cache"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/redis/go-redis/v9"
	"time"
)

// ========================================================================================
// 										Structure
// ========================================================================================

// New func for connecting to Redis server.
func New() *Driver {
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

	return &Driver{
		redisCache: redis.NewClient(options),
	}
}

type Driver struct {
	redisCache *redis.Client
}

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
		log.Warnf("Error while reading key `%v`", key)
		return nil, err
	}

	return val, nil
}

func (r *Driver) Del(key string) error {
	if err := r.redisCache.Del(context.Background(), cache.Key(key)).Err(); err != nil {
		log.Warnf("Error while deleting key `%v`", key)
		return err
	}

	return nil
}
