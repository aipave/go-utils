package test

import (
	"testing"

	"github.com/alyu01/go-utils/gcache"
)

func TestGcache(t *testing.T) {
	gcache.Set("foo", "bar", 5)
	v, ok := gcache.Get("foo")
	t.Log(v, ok)

}
