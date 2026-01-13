Redis Watcher
====

[![Go Report Card](https://goreportcard.com/badge/github.com/casbin/redis-watcher)](https://goreportcard.com/report/github.com/casbin/redis-watcher)
[![Go](https://github.com/casbin/redis-watcher/actions/workflows/ci.yml/badge.svg)](https://github.com/casbin/redis-watcher/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/casbin/redis-watcher/badge.svg?branch=master)](https://coveralls.io/github/casbin/redis-watcher?branch=master)
[![Godoc](https://godoc.org/github.com/casbin/redis-watcher?status.svg)](https://godoc.org/github.com/casbin/redis-watcher)
[![Release](https://img.shields.io/github/release/casbin/redis-watcher.svg)](https://github.com/casbin/redis-watcher/releases/latest)
[![Discord](https://img.shields.io/discord/1022748306096537660?logo=discord&label=discord&color=5865F2)](https://discord.gg/S5UjpzGZjN)
[![Sourcegraph](https://sourcegraph.com/github.com/casbin/redis-watcher/-/badge.svg)](https://sourcegraph.com/github.com/casbin/redis-watcher?badge)

Redis Watcher is a [Redis](http://redis.io) watcher for [Casbin](https://github.com/casbin/casbin).

## Installation

```bash
go get github.com/casbin/redis-watcher/v2
```

## Simple Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
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

	// Or initialize the watcher in redis cluster.
	// w, _ := rediswatcher.NewWatcherWithCluster("localhost:6379,localhost:6379,localhost:6379", rediswatcher.WatcherOptions{
	// 	ClusterOptions: redis.ClusterOptions{
	// 		Password: "",
	// 	},
	// 	Channel: "/casbin",
	// 	IgnoreSelf: false,
	// })

	// Initialize the enforcer.
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	// Set the watcher for the enforcer.
	_ = e.SetWatcher(w)

	// Set callback to local example
	_ = w.SetUpdateCallback(updateCallback)
	
	// Or use the default callback
	// _ = w.SetUpdateCallback(rediswatcher.DefaultUpdateCallback(e))

	// Update the policy to test the effect.
	// You should see "[casbin rules updated]" in the log.
	_ = e.SavePolicy()
	// Only exists in test
	fmt.Scanln()
}

```

## Using Existing Redis Client

If you already have an existing Redis client instance (e.g., due to restricted access to connection details or centralized Redis setup), you can reuse it to create a watcher.

### With Regular Redis Client

```go
package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Create your own Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Initialize the watcher with existing client
	// Pass empty string for addr since the SubClient and PubClient are already provided
	w, _ := rediswatcher.NewWatcher("", rediswatcher.WatcherOptions{
		SubClient: redisClient,
		PubClient: redisClient,
		Channel:   "/casbin",
	})

	// Initialize the enforcer.
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	// Set the watcher for the enforcer.
	_ = e.SetWatcher(w)

	// Set callback
	_ = w.SetUpdateCallback(rediswatcher.DefaultUpdateCallback(e))

	// Update the policy to test the effect.
	_ = e.SavePolicy()
	fmt.Scanln()
}
```

### With Redis Cluster Client

```go
package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Create your own Redis cluster client
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"localhost:7000", "localhost:7001", "localhost:7002"},
		Password: "",
	})

	// Initialize the watcher with existing cluster client
	// Pass empty string for addrs since the SubClient and PubClient are already provided
	w, _ := rediswatcher.NewWatcherWithCluster("", rediswatcher.WatcherOptions{
		SubClient: clusterClient,
		PubClient: clusterClient,
		Channel:   "/casbin",
	})

	// Initialize the enforcer.
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	// Set the watcher for the enforcer.
	_ = e.SetWatcher(w)

	// Set callback
	_ = w.SetUpdateCallback(rediswatcher.DefaultUpdateCallback(e))

	// Update the policy to test the effect.
	_ = e.SavePolicy()
	fmt.Scanln()
}
```

## Getting Help

- [Casbin](https://github.com/casbin/casbin)
- [go-redis](https://github.com/go-redis/redis)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
