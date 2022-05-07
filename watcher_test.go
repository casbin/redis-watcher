package rediswatcher

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
)

func initWatcherWithOptions(t *testing.T, wo WatcherOptions) (*casbin.Enforcer, *Watcher) {
	w, err := NewWatcher("127.0.0.1:6379", wo)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	e, err := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")
	if err != nil {
		t.Fatalf("Failed to create enforcer: %v", err)
	}
	_ = e.SetWatcher(w)
	_ = w.SetUpdateCallback(MakeDefaultUpdateCallback(e))
	return e, w.(*Watcher)
}

func initWatcher(t *testing.T) (*casbin.Enforcer, *Watcher) {
	return initWatcherWithOptions(t, WatcherOptions{})
}

func TestWatcher(t *testing.T) {
	_, w := initWatcher(t)
	_ = w.SetUpdateCallback(func(s string) {
		fmt.Println(s)
	})
	_ = w.Update()
	w.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestWatcherWithIgnoreSelfTrue(t *testing.T) {
	wo := WatcherOptions{
		IgnoreSelf: true,
		OptionalUpdateCallback: func(s string) {
			t.Fatalf("This callback should not be called when IgnoreSelf is set true.")
		},
	}
	_, w := initWatcherWithOptions(t, wo)
	_ = w.Update()
	w.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdate(t *testing.T) {
	_, w := initWatcher(t)
	_ = w.SetUpdateCallback(func(s string) {
		CustomDefaultFunc(
			func(id string, params interface{}) {
				t.Fatalf("method mapping error")
			},
		)(s, func(ID string, params interface{}) {
			if ID != w.options.LocalID {
				t.Fatalf("instance ID should be %s instead of %s", w.options.LocalID, ID)
			}
		}, nil, nil, nil, nil)
	})
	_ = w.Update()
	w.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForAddPolicy(t *testing.T) {
	wo := WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.AddPolicy("alice", "book1", "write")
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
}

func TestUpdateForRemovePolicy(t *testing.T) {
	wo := WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.RemovePolicy("alice", "data1", "read")
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
}

func TestUpdateForRemoveFilteredPolicy(t *testing.T) {
	wo := WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.RemoveFilteredPolicy(1, "data1", "read")
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
}

func TestUpdateSavePolicy(t *testing.T) {
	e, w := initWatcher(t)
	_ = w.SetUpdateCallback(func(s string) {
		CustomDefaultFunc(
			func(id string, params interface{}) {
				t.Fatalf("method mapping error")
			},
		)(s, nil, nil, nil, nil, func(ID string, params interface{}) {
			if ID != w.options.LocalID {
				t.Fatalf("instance ID should be %s instead of %s", w.options.LocalID, ID)
			}
			res := params.(string)
			if res != "" {
				t.Fatalf("instance Params should be empty instead of %s", res)
			}
		})
	})
	_ = e.SavePolicy()
	w.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForAddPolicies(t *testing.T) {
	rules := [][]string{
		{"jack", "data4", "read"},
		{"katy", "data4", "write"},
		{"leyo", "data4", "read"},
		{"ham", "data4", "write"},
	}

	wo := WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.AddPolicies(rules)
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
}

func TestUpdateForRemovePolicies(t *testing.T) {
	rules := [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
	}
	wo := WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.RemovePolicies(rules)
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
}
