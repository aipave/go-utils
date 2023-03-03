package test

import (
    "testing"

    gcache "github.com/alyu01/go-utils/gcache"
)

func TestGcache(t *testing.T) {
    gcache.Set("foo", "bar", 5)
    v, ok := gcache.Get("foo")
    t.Log(v, ok)

}
