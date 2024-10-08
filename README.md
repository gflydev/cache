# gFly Caching

    Copyright © 2023, gFly
    https://www.gFly.dev
    All rights reserved.

### Usage

Install
```bash
go get -u github.com/gflydev/cache@v1.0.1
```

Quick usage `main.go`
```go
import (
    cacheRedis "github.com/gflydev/cache/redis"
    "github.com/gflydev/cache"
)

// Register Redis cache
cache.Register(cacheRedis.New())

// Set cache 15 days
if err = cache.Set(key, value, time.Duration(15*24*3600) * time.Second); err != nil {
    log.Errorf("Error %q", err)
}

val, err := cache.Get(key)
if err != nil {
    log.Errorf("Error %q", err)
}

// Delete 
if err = cache.Del(key); err != nil {
    log.Errorf("Error %q", err)
}
```