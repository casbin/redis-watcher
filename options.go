package rediswatcher

import (
	rds "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type WatcherOptions struct {
	rds.Options
	SubClient              *rds.Client
	PubClient              *rds.Client
	Channel                string
	IgnoreSelf             bool
	LocalID                string
	OptionalUpdateCallback func(string)
}

func initConfig(option *WatcherOptions) {
	if option.LocalID == "" {
		option.LocalID = uuid.New().String()
	}
	if option.Channel == "" {
		option.Channel = "/casbin"
	}
}
