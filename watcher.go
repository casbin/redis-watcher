package rediswatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2/model"

	"github.com/casbin/casbin/v2/persist"
	rds "github.com/go-redis/redis/v8"
)

type Watcher struct {
	l         sync.Mutex
	subClient *rds.Client
	pubClient *rds.Client
	options   WatcherOptions
	close     chan struct{}
	callback  func(string)
	ctx       context.Context
}

type MSG struct {
	Method string
	ID     string
	Sec    string
	Ptype  string
	Params interface{}
}

func (m *MSG) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary decodes the struct into a User
func (m *MSG) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	return nil
}

// NewWatcher creates a new Watcher to be used with a Casbin enforcer
// addr is a redis target string in the format "host:port"
// setters allows for inline WatcherOptions
//
// 		Example:
// 				w, err := rediswatcher.NewWatcher("127.0.0.1:6379",WatcherOptions{})
//
func NewWatcher(addr string, option WatcherOptions) (persist.Watcher, error) {
	option.Addr = addr
	w := &Watcher{
		subClient: rds.NewClient(&option.Options),
		pubClient: rds.NewClient(&option.Options),
		ctx:       context.Background(),
		close:     make(chan struct{}),
	}

	if err := w.subClient.Ping(w.ctx).Err(); err != nil {
		return nil, err
	}
	if err := w.pubClient.Ping(w.ctx).Err(); err != nil {
		return nil, err
	}

	initConfig(&option)
	w.options = option

	w.subscribe()

	return w, nil
}

// NewPublishWatcher return a Watcher only publish but not subscribe
func NewPublishWatcher(addr string, option WatcherOptions) (persist.Watcher, error) {
	option.Addr = addr
	w := &Watcher{
		pubClient: rds.NewClient(&option.Options),
		ctx:       context.Background(),
		close:     make(chan struct{}),
	}

	initConfig(&option)
	w.options = option

	return w, nil
}

// SetUpdateCallback SetUpdateCallBack sets the update callback function invoked by the watcher
// when the policy is updated. Defaults to Enforcer.LoadPolicy()
func (w *Watcher) SetUpdateCallback(callback func(string)) error {
	w.l.Lock()
	w.callback = callback
	w.l.Unlock()
	return nil
}

// Update publishes a message to all other casbin instances telling them to
// invoke their update callback
func (w *Watcher) Update() error {
	return w.logRecord(func() error {
		w.l.Lock()
		defer w.l.Unlock()
		return w.pubClient.Publish(context.Background(), w.options.Channel, &MSG{"Update", w.options.LocalID, "", "", ""}).Err()
	})
}

// UpdateForAddPolicy calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.AddPolicy()
func (w *Watcher) UpdateForAddPolicy(sec, ptype string, params ...string) error {
	return w.logRecord(func() error {
		w.l.Lock()
		defer w.l.Unlock()
		return w.pubClient.Publish(context.Background(), w.options.Channel, &MSG{"UpdateForAddPolicy", w.options.LocalID, sec, ptype, params}).Err()
	})
}

// UpdateForRemovePolicy UPdateForRemovePolicy calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.RemovePolicy()
func (w *Watcher) UpdateForRemovePolicy(sec, ptype string, params ...string) error {
	return w.logRecord(func() error {
		w.l.Lock()
		defer w.l.Unlock()
		return w.pubClient.Publish(context.Background(), w.options.Channel, &MSG{"UpdateForRemovePolicy", w.options.LocalID, sec, ptype, params}).Err()
	})
}

// UpdateForRemoveFilteredPolicy calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.RemoveFilteredNamedGroupingPolicy()
func (w *Watcher) UpdateForRemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error {
	return w.logRecord(func() error {
		w.l.Lock()
		defer w.l.Unlock()
		return w.pubClient.Publish(context.Background(), w.options.Channel,
			&MSG{"UpdateForRemoveFilteredPolicy", w.options.LocalID,
				sec,
				ptype,
				fmt.Sprintf("%d %s", fieldIndex, strings.Join(fieldValues, " ")),
			},
		).Err()
	})
}

// UpdateForSavePolicy calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.RemoveFilteredNamedGroupingPolicy()
func (w *Watcher) UpdateForSavePolicy(model model.Model) error {
	return w.logRecord(func() error {
		w.l.Lock()
		defer w.l.Unlock()
		return w.pubClient.Publish(context.Background(), w.options.Channel, &MSG{"UpdateForSavePolicy", w.options.LocalID, "", "", model}).Err()
	})
}

func (w *Watcher) logRecord(f func() error) error {
	err := f()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (w *Watcher) unsubscribe(psc *rds.PubSub) error {
	return psc.Unsubscribe(w.ctx)
}

func (w *Watcher) subscribe() {
	w.l.Lock()
	sub := w.subClient.Subscribe(w.ctx, w.options.Channel)
	w.l.Unlock()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			err := sub.Close()
			if err != nil {
				log.Println(err)
			}
			err = w.pubClient.Close()
			if err != nil {
				log.Println(err)
			}
			err = w.subClient.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		ch := sub.Channel()
		wg.Done()
		for msg := range ch {
			select {
			case <-w.close:
				return
			default:
			}
			data := msg.Payload
			w.callback(data)
		}
	}()
	wg.Wait()
}

func (w *Watcher) GetWatcherOptions() WatcherOptions {
	w.l.Lock()
	defer w.l.Unlock()
	return w.options
}

func (w *Watcher) Close() {
	w.l.Lock()
	defer w.l.Unlock()
	close(w.close)
	w.pubClient.Publish(w.ctx, w.options.Channel, "Close")
}
