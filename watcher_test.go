package rediswatcher_test

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	rediswatcher "github.com/casbin/redis-watcher/v2"
)

func initWatcherWithOptions(t *testing.T, wo rediswatcher.WatcherOptions, cluster ...bool) (*casbin.Enforcer, *rediswatcher.Watcher) {
	var (
		w   persist.Watcher
		err error
	)
	if len(cluster) > 0 && cluster[0] {
		w, err = rediswatcher.NewWatcherWithCluster("127.0.0.1:6379,127.0.0.1:6379,127.0.0.1:6379", wo)
	} else {
		w, err = rediswatcher.NewWatcher("127.0.0.1:6379", wo)
	}
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	e, err := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")
	if err != nil {
		t.Fatalf("Failed to create enforcer: %v", err)
	}
	_ = e.SetWatcher(w)
	_ = w.SetUpdateCallback(func(s string) {
		t.Log(s)
		rediswatcher.DefaultUpdateCallback(e)(s)
	})
	return e, w.(*rediswatcher.Watcher)
}

func initWatcher(t *testing.T, cluster ...bool) (*casbin.Enforcer, *rediswatcher.Watcher) {
	return initWatcherWithOptions(t, rediswatcher.WatcherOptions{}, cluster...)
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

func TestWatcherWithIngnoreSelfFalse(t *testing.T) {
	wo := rediswatcher.WatcherOptions{
		LocalID: uuid.New().String(),
	}
	_, w := initWatcherWithOptions(t, wo)
	_ = w.SetUpdateCallback(func(s string) {
		msg := &rediswatcher.MSG{}
		_ = json.Unmarshal([]byte(s), msg)
		if msg.ID != wo.LocalID {
			t.Fatalf("instance ID should be %s instead of %s", wo.LocalID, msg.ID)
		}
	})

	_ = w.Update()
	time.Sleep(time.Millisecond * 500)
	w.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestWatcherWithIgnoreSelfTrue(t *testing.T) {
	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	_, w := initWatcherWithOptions(t, wo)

	_ = w.SetUpdateCallback(func(s string) {
		t.Fatalf("This callback should not be called when IgnoreSelf is set true.")
	})

	_ = w.Update()
	time.Sleep(time.Millisecond * 500)
	w.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdate(t *testing.T) {
	_, w := initWatcher(t)
	_ = w.SetUpdateCallback(func(s string) {
		msgStruct := &rediswatcher.MSG{}

		err := msgStruct.UnmarshalBinary([]byte(s))
		if err != nil {
			t.Error(err)
			return
		}
		if msgStruct.Method != "Update" {
			t.Errorf("Method should be Update instead of %s", msgStruct.Method)
		}
	})
	_ = w.Update()
	w.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForAddPolicy(t *testing.T) {
	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.AddPolicy("alice", "book1", "write")
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Log("Method", "AddPolicy")
		t.Log("e.policy", e.GetPolicy())
		t.Log("e2.policy", e2.GetPolicy())
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForRemovePolicy(t *testing.T) {
	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.RemovePolicy("alice", "data1", "read")
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Log("Method", "RemovePolicy")
		t.Log("e.policy", e.GetPolicy())
		t.Log("e2.policy", e2.GetPolicy())
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForRemoveFilteredPolicy(t *testing.T) {
	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.RemoveFilteredPolicy(1, "data1", "read")
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Log("Method", "RemoveFilteredPolicy")
		t.Log("e.policy", e.GetPolicy())
		t.Log("e2.policy", e2.GetPolicy())
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateSavePolicy(t *testing.T) {
	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_ = e.SavePolicy()
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Log("Method", "SavePolicy")
		t.Log("e.policy", e.GetPolicy())
		t.Log("e2.policy", e2.GetPolicy())
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForAddPolicies(t *testing.T) {
	rules := [][]string{
		{"jack", "data4", "read"},
		{"katy", "data4", "write"},
		{"leyo", "data4", "read"},
		{"ham", "data4", "write"},
	}

	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.AddPolicies(rules)
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Log("Method", "AddPolicies")
		t.Log("e.policy", e.GetPolicy())
		t.Log("e2.policy", e2.GetPolicy())
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForRemovePolicies(t *testing.T) {
	rules := [][]string{
		{"jack", "data4", "read"},
		{"katy", "data4", "write"},
		{"leyo", "data4", "read"},
		{"ham", "data4", "write"},
	}

	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.RemovePolicies(rules)
	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Log("Method", "RemovePolicies")
		t.Log("e.policy", e.GetPolicy())
		t.Log("e2.policy", e2.GetPolicy())
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForUpdatePolicy(t *testing.T) {
	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.UpdatePolicy([]string{"alice", "data1", "read"}, []string{"alice", "data1", "write"})

	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Log("Method", "UpdatePolicy")
		t.Log("e.policy", e.GetPolicy())
		t.Log("e2.policy", e2.GetPolicy())
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
	time.Sleep(time.Millisecond * 500)
}

func TestUpdateForUpdatePolicies(t *testing.T) {
	wo := rediswatcher.WatcherOptions{
		IgnoreSelf: true,
	}
	e, w := initWatcherWithOptions(t, wo)
	e2, w2 := initWatcherWithOptions(t, wo)

	time.Sleep(time.Millisecond * 500)
	_, _ = e.UpdatePolicies([][]string{{"alice", "data1", "read"}}, [][]string{{"alice", "data1", "write"}})

	time.Sleep(time.Millisecond * 500)
	if !reflect.DeepEqual(e2.GetPolicy(), e.GetPolicy()) {
		t.Log("Method", "UpdatePolicies")
		t.Log("e.policy", e.GetPolicy())
		t.Log("e2.policy", e2.GetPolicy())
		t.Error("These two enforcers' policies should be equal")
	}

	w.Close()
	w2.Close()
	time.Sleep(time.Millisecond * 500)
}
