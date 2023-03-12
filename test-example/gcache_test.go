package test_example

import (
	"testing"
	"time"

	gcache "github.com/aipave/go-utils/gcache"
)

func TestGcache(t *testing.T) {
	gcache.Set("foo", "bar", 5)
	v, ok := gcache.Get("foo")
	t.Log(v, ok)

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
		gcache.Set(mock.key, mock.val, mock.ttl)
	}

	time.Sleep(1)

	for _, mock := range mocks {
		val, found := gcache.Get(mock.key)
		if !found || val.(string) != mock.val {
			t.Fatalf("Unexpected value: %v (%v) to key: %v", val, found, mock.key)
		}
	}

}
