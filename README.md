Redis Watcher 
---

> For Go 1.17+, use v2.4.0+ ,  
> For Go 1.16 and below, stay on v2.3.0.  

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

	// ======================================
	// Example 2: Redis Cluster (Using addrs)
	// ======================================
	// Or initialize the watcher in redis cluster.
	// w, _ := rediswatcher.NewWatcherWithCluster("localhost:6379,localhost:6379,localhost:6379", rediswatcher.WatcherOptions{
	// 	ClusterOptions: redis.ClusterOptions{
	// 		Password: "",
	// 	},
	// 	Channel: "/casbin",
	// 	IgnoreSelf: false,
	// })


	// ======================================
	// Example 3: Standalone Redis with existing Redis client
	// ======================================
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0,
	// })
	// here we can directly pass existing redis client as SubClient and/or PubClient
	// w, _ := rediswatcher.NewWatcher("", rediswatcher.WatcherOptions{
	// 	SubClient: rdb,
	// 	PubClient: rdb,
	// 	Channel:   "/casbin",
	// })

	// ======================================
	// Example 4: Redis Cluster with existing cluster clients
	// ======================================
	// clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs:    []string{"localhost:6379", "localhost:6380", "localhost:6381"},
	// 	Password: "",
	// })
	// here we can directly pass existing redis client as SubClusterClient and/or PubClusterClient
	// w, _ := rediswatcher.NewWatcherWithCluster("", rediswatcher.WatcherOptions{
	// 	SubClusterClient: clusterClient,
	// 	PubClusterClient: clusterClient,
	// 	Channel:          "/casbin",
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

## Getting Help

- [Casbin](https://github.com/casbin/casbin)
- [go-redis](https://github.com/go-redis/redis)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
