package rediswatcher

import (
	"errors"
	"github.com/casbin/casbin/v2"
	rds "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
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

func MakeDefaultUpdateCallback(e *casbin.Enforcer) func(string) {
	return func(msg string) {
		msgStruct := &MSG{}

		err := msgStruct.UnmarshalBinary([]byte(msg))
		if err != nil {
			log.Println(err)
		}

		e.EnableAutoNotifyWatcher(false)

		switch msgStruct.Method {
		case UpdateType_Update:
		case UpdateType_UpdateForSavePolicy:
			err = e.LoadPolicy()
		case UpdateType_UpdateForAddPolicy:
			params := msgStruct.Params[0]

			_, err = e.AddPolicy(params)
		case UpdateType_UpdateForRemovePolicy:
			params := msgStruct.Params[0]

			_, err = e.RemovePolicy(params)
		case UpdateType_UpdateForRemoveFilteredPolicy:
			params := msgStruct.Params[0][0]

			// parse the result of fmt.Sprintf("%d %s", fieldIndex, strings.Join(fieldValues, " "))
			paramsList := strings.Fields(params)
			fieldIndex, _ := strconv.Atoi(paramsList[0])
			fieldValues := paramsList[1:]

			_, err = e.RemoveFilteredPolicy(fieldIndex, fieldValues...)
		case UpdateType_UpdateForAddPolicies:
			params := msgStruct.Params

			_, err = e.AddPolicies(params)
		case UpdateType_UpdateForRemovePolicies:
			params := msgStruct.Params

			_, err = e.RemovePolicies(params)
		default:
			err = errors.New("unknown update type")
		}
		if err != nil {
			log.Println(err)
		}

		e.EnableAutoNotifyWatcher(true)
	}
}
