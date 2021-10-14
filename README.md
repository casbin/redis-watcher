Redis Watcher 
---
[![Go](https://github.com/casbin/redis-watcher/actions/workflows/ci.yml/badge.svg)](https://github.com/casbin/redis-watcher/actions/workflows/ci.yml)
[![report](https://goreportcard.com/badge/github.com/casbin/redis-watcher)](https://goreportcard.com/report/github.com/casbin/redis-watcher)
[![Coverage Status](https://coveralls.io/repos/github/casbin/redis-watcher/badge.svg?branch=master)](https://coveralls.io/github/casbin/redis-watcher?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/casbin/redis-watcher/v2.svg)](https://pkg.go.dev/github.com/casbin/redis-watcher/v2)
[![Release](https://img.shields.io/github/v/release/casbin/redis-watcher)](https://github.com/casbin/redis-watcher/releases/latest)

Redis Watcher is a [Redis](http://redis.io) watcher for [Casbin](https://github.com/casbin/casbin).

## Installation

    go get github.com/casbin/redis-watcher/v2

## Simple Example

```go
package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	watcher "github.com/casbin/redis-watcher/v2"
	"github.com/go-redis/redis/v8"
)

func updateCallback(msg string) {
	log.Println(msg)
}

func main() {
	// Initialize the watcher.
	// Use the Redis host as parameter.
	w, _ := rediswatcher.NewWatcher("localhost:6379", rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: "",
		},
		Channel:    "/casbin",
		// Only exists in test, generally be true
		IgnoreSelf: false,
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
	// Only exists in test
	fmt.Scanln()
}

```

## Getting Help

- [Casbin](https://github.com/casbin/casbin)
- [go-redis](https://github.com/go-redis/redis)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
>>>>>>> 243bd42 (refactor)
