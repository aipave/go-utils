package gcache

import (
    "context"
)

var defaultLocalCache = New[string, any]()

// Get
func Get(key string) (interface{}, bool) {
    return defaultLocalCache.Get(key)
}

// GetWithLoad
func GetWithLoad(ctx context.Context, key string, load LoadFunc[string, any]) (interface{}, error) {
    value, found := defaultLocalCache.Get(key)
    if found {
        return value, nil
    }

    c, _ := defaultLocalCache.(*cache[string, any])

    return c.loadData(ctx, key, load)
}

// Set key, value, ttl(s)
func Set(key string, value interface{}, ttl int64) bool {
    return defaultLocalCache.SetWithExpire(key, value, ttl)
}

// Del key
func Del(key string) {
    defaultLocalCache.Del(key)
}

// Clear 清空所有队列和缓存
func Clear() {
    defaultLocalCache.Clear()
}
