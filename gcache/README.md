# gcache

[![BK Pipelines Status](https://api.bkdevops.qq.com/process/api/external/pipelines/projects/pcgtrpcproject/p-c9b0d4b7b7754407ae09b4c8887a0ba6/badge?X-DEVOPS-PROJECT-ID=pcgtrpcproject)](http://devops.oa.com:/ms/process/api-html/user/builds/projects/pcgtrpcproject/pipelines/p-c9b0d4b7b7754407ae09b4c8887a0ba6/latestFinished?X-DEVOPS-PROJECT-ID=pcgtrpcproject)[![Coverage](https://tcoverage.woa.com/api/getCoverage/getTotalImg/?pipeline_id=p-c9b0d4b7b7754407ae09b4c8887a0ba6)](http://macaron.oa.com/api/coverage/getTotalLink/?pipeline_id=p-c9b0d4b7b7754407ae09b4c8887a0ba6)[![GoDoc](https://img.shields.io/badge/API%20Docs-GoDoc-green)](http://godoc.oa.com/git.code.oa.com/trpc-go/trpc-database/cos)

localcache is a local K-V caching component that runs on a single machine. 
It allows concurrent access by multiple goroutines and supports LRU-based and expiration-based eviction policies.
When the localcache capacity reaches its limit, data eviction is performed based on LRU, and expired key-value pairs are deleted using a time wheel implementation.

## how to use


```go
package main

import (
    "xxxxx/gcache"
)

func LoadData(ctx context.Context, key string) (interface{}, error) {
    return "cat", nil
}

func main() {
    // set 5s tll
    gcache.Set("foo", "bar", 5)

    // get value of key
    value, found := gcache.Get("foo")

    // get value of key
    // set 5s tll
    value, err := gcache.GetWithLoad(context.TODO(), "tom", LoadData, 5)

    // del key
    gcache.Del("foo")

    // clear
    gcache.Clear()
}
```

## gomod import

	require xxxxx/xxxx/gcache v3.0.0

## use config

New () to generate Cache instance, call the instance function function

### Optional parameters

#### **WithCapacity(capacity int)**

set cache cap. min = 1，max=1e30
will del tail elem. when to max, default value is 1e30

#### **WithExpiration(ttl int64)**

#### **WithLoad(f LoadFunc)**

```go
type LoadFunc func(ctx context.Context, key string) (interface{}, error)
```

Set the data load function, the key does not exist in the cache, use this function to load the corresponding value, and cached in the cache.

#### Cache interface

```go
type Cache interface {
	// Get 
	Get(key string) (interface{}, bool)

	// GetWithLoad 
	GetWithLoad(ctx context.Context, key string) (interface{}, error)

	// Set 
	Set(key string, value interface{}) bool

	// SetWithExpire 
	SetWithExpire(key string, value interface{}, ttl int64) bool

	// Del key
	Del(key string)

	// Clear 
	Clear()

	// Close close cache
	Close()
}
```

## example

#### set capacity and ttl

```go
package main

import (
    "fmt"
    "time"

    "git.code.oa.com/trpc-go/trpc-database/gcache"
)

func main() {
    var lc gcache.Cache

    // new cap = 100, ttl = 5s
    lc = gcache.New(gcache.WithCapacity(100), gcache.WithExpiration(5))

    // set key-value, ttl = 5s
    lc.Set("foo", "bar")

    // re-custom tll for key, 10s
    lc.SetWithExpire("tom", "cat", 10)

    // asyn wait 
    time.Sleep(time.Millisecond)

    // get value
    val, found := lc.Get("foo")

    // del key: "foo"
    lc.Del("foo")
}
```

#### Custom loading function

Set a custom data loading function and use GetWithLoad(key) to obtain the value.

```go
func main() {
loadFunc := func (ctx context.Context, key string) (interface{}, error) {
return "cat", nil
}

lc := gcache.New(gcache.WithLoad(loadFunc), gcache.WithExpiration(5))

// return err for loadFunc
val, err := lc.GetWithLoad(context.TODO(), "tom")
}
```

### TODO

- [ ] Add Metrics data statistics
- [ ] Add control of memory usage
- [ ] Introduce an upgraded version of LRU: W-tinyLRU, which more efficiently controls the write and eviction of keys.
  
  
