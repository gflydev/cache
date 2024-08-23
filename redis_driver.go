package cache

import (
	"context"
	"fmt"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisCache *redis.Client

// New func for connecting to Redis server.
func init() {
	// Define Redis database number.
	dbNumber := utils.Getenv("REDIS_DB_NUMBER", 0)

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

	redisCache = redis.NewClient(options)
}

func New() ICache {
	return RedisDriver{}
}

type RedisDriver struct{}

func (r RedisDriver) Set(key string, value interface{}, expiration time.Duration) error {
	if err := redisCache.Set(context.Background(), Key(key), value, expiration).Err(); err != nil {
		log.Errorf("Error while writing Redis cache %q", err)
		return err
	}

	return nil
}
func (r RedisDriver) Get(key string) (interface{}, error) {
	val, err := redisCache.Get(context.Background(), Key(key)).Result()
	if err != nil {
		log.Errorf("Error while reading Redis cache %q", err)
		return nil, err
	}

	return val, nil
}

func (r RedisDriver) Del(key string) error {
	if err := redisCache.Del(context.Background(), Key(key)).Err(); err != nil {
		log.Errorf("Error while deleting Redis cache %q", err)
		return err
	}

	return nil
}
