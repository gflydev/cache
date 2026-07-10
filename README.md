# gFly Caching

    Copyright © 2023, gFly
    https://www.gFly.dev
    All rights reserved.

A small, driver-based caching abstraction for the [gFly](https://www.gFly.dev)
framework. The `cache` package exposes a backend-agnostic API; concrete
backends (currently Redis) live in sub-packages and are wired up with
`cache.Register`.

### Install

```bash
go get -u github.com/gflydev/cache
```

### Usage

```go
import (
    "errors"
    "time"

    "github.com/gflydev/cache"
    cacheRedis "github.com/gflydev/cache/redis"
)

// Register the Redis driver once during start-up.
cache.Register(cacheRedis.New())

// Set a value that expires in 15 days.
if err := cache.Set(key, value, 15*24*time.Hour); err != nil {
    log.Errorf("Error %q", err)
}

// Get a value, distinguishing a genuine miss from a backend failure.
val, err := cache.Get(key)
switch {
case errors.Is(err, cache.ErrCacheMiss):
    // key not present – recompute and Set it
case err != nil:
    log.Errorf("Error %q", err)
default:
    _ = val
}

// Delete a key.
if err := cache.Del(key); err != nil {
    log.Errorf("Error %q", err)
}

// Release backend resources on shutdown.
defer cache.Close()
```

### Configuration

The Redis driver reads the following environment variables (see `.env.example`):

| Variable            | Default     | Description                 |
|---------------------|-------------|-----------------------------|
| `APP_CODE`          | `gfly`      | Key namespace prefix        |
| `REDIS_HOST`        | `localhost` | Redis host                  |
| `REDIS_PORT`        | `6379`      | Redis port                  |
| `REDIS_PASSWORD`    | *(empty)*   | Redis password              |
| `REDIS_DEFAULT_DB`  | `0`         | Redis database number       |

Connection tuning is available via functional options:

```go
cacheRedis.New(
    cacheRedis.WithPoolSize(50),
    cacheRedis.WithTimeouts(5*time.Second, 3*time.Second, 3*time.Second),
)
```

Verify connectivity at start-up with `driver.Ping()`.

### Errors

| Error                   | Meaning                                          |
|-------------------------|--------------------------------------------------|
| `cache.ErrCacheMiss`    | Requested key does not exist.                    |
| `cache.ErrNotRegistered`| A helper was called before `cache.Register`.     |

### Testing

```bash
go test ./...
```
