# gFly Caching

    Copyright © 2023, gFly
    https://www.gFly.dev
    All rights reserved.

### Usage

Install
```bash
go get -u github.com/gflydev/cache@v1.0.0
```

Quick usage `main.go`
```go
import (
    _ "github.com/gflydev/cache"
)

kv := cache.New()

// Set cache 15 days
if err = kv.Set(key, value, time.Duration(15*24*3600) * time.Second); err != nil {
    log.Errorf("Error %q", err)
}

val, err := kv.Get(key)
if err != nil {
    log.Errorf("Error %q", err)
}

// Delete 
if err = kv.Del(key); err != nil {
    log.Errorf("Error %q", err)
}
```