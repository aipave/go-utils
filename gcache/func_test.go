package gcache

import (
    "context"
    "errors"
    "testing"
    "time"
)

// TestFuncGetAndSet 测试Get和Set
func TestFuncGetAndSet(t *testing.T) {
    mocks := []struct {
        key string
        val string
        ttl int64
    }{
        {"A", "a", 1000},
        {"B", "b", 1000},
        {"", "null", 1000},
    }

    for _, mock := range mocks {
        Set(mock.key, mock.val, mock.ttl)
    }

    time.Sleep(wait)

    for _, mock := range mocks {
        val, found := Get(mock.key)
        if !found || val.(string) != mock.val {
            t.Fatalf("Unexpected value: %v (%v) to key: %v", val, found, mock.key)
        }
    }
}

// TestFuncExpireSet 测试带Expire的Set
func TestFuncExpireSet(t *testing.T) {
    Set("Foo", "Bar", 1)
    time.Sleep(1)
    if val, found := Get("Foo"); found {
        t.Fatalf("unexpected expired value: %v to key: %v", val, "Foo")
    }
}

// TestFuncGetWithLoad 测试WithLoad情况下的Get
func TestFuncGetWithLoad(t *testing.T) {
    m := map[string]interface{}{
        "A1": "a",
        "B":  "b",
        "":   "null",
    }
    loadFunc := func(ctx context.Context, key string) (interface{}, error) {
        if v, exist := m[key]; exist {
            return v, nil
        }
        return nil, errors.New("key not exist")
    }
    value, err := GetWithLoad(context.TODO(), "A1", loadFunc)
    if err != nil || value.(string) != "a" {
        t.Fatalf("unexpected GetWithLoad value: %v, want:a, err:%v", value, err)
    }

    time.Sleep(wait)

    got2, found := Get("A1")
    if !found || got2.(string) != "a" {
        t.Fatalf("unexpected Get value: %v, want:a, found:%v", got2, found)
    }
}

func TestDelAndClear(t *testing.T) {
    Convey("DelAndClear", t, func() {
        Set("Foo", "bar", 10)
        Set("Foo1", "bar", 10)
        time.Sleep(time.Millisecond * 10)
        _, ok := Get("Foo")
        So(ok, ShouldBeTrue)
        Del("Foo")
        time.Sleep(time.Millisecond * 10)
        _, ok = Get("Foo")
        So(ok, ShouldBeFalse)
        Clear()
        time.Sleep(time.Millisecond * 10)
        _, ok = Get("Foo1")
        So(ok, ShouldBeFalse)
    })
}
