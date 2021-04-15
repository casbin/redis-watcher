Redis Watcher [![Go](https://github.com/casbin/redis-watcher/actions/workflows/ci.yml/badge.svg)](https://github.com/casbin/redis-watcher/actions/workflows/ci.yml)

Redis Watcher is a [Redis](http://redis.io) watcher for [Casbin](https://github.com/casbin/casbin).

## Installation

    go get github.com/casbin/redis-watcher

## Simple Example

```go
package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	watcher "github.com/casbin/redis-watcher"
	"github.com/go-redis/redis/v8"
)

func updateCallback(msg string) {
	log.Println(msg)
}

func main() {
	// Initialize the watcher.
	// Use the Redis host as parameter.
	w, _ := watcher.NewWatcher("localhost:6379", watcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: "",
		},
		Channel:    "/casbin",
		IgnoreSelf: true,
	})

	// Initialize the enforcer.
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	// Set the watcher for the enforcer.
	_ = e.SetWatcher(w)

	// Set callback to local example
	_ = w.SetUpdateCallback(updateCallback)

	// Update the policy to test the effect.
	// You should see "[casbin rules updated]" in the log.
	_ = e.SavePolicy()
}

```

## Getting Help

- [Casbin](https://github.com/casbin/casbin)
- [go-redis](https://github.com/go-redis/redis)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
>>>>>>> 243bd42 (refactor)
