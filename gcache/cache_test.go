package gcache

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "reflect"
    "sync"
    "testing"
    "time"
)

var wait = time.Millisecond

// TestCacheSetRace 测试Cache Set并发竞争
func TestCacheSetRace(t *testing.T) {
    cache := New[string, string](WithExpiration[string, string](100))
    n := 8128
    var wg sync.WaitGroup
    wg.Add(n)
    for i := 0; i < n; i++ {
        go func() {
            cache.Set("foo", "bar")
            cache.Get("foo")
            wg.Done()
        }()
    }
    wg.Wait()
}

// TestCacheSetGet 测试Cache 先Set后Get
func TestCacheSetGet(t *testing.T) {
    mocks := []struct {
        key string
        val string
        ttl int64
    }{
        {"A", "a", 1000},
        {"B", "b", 1000},
        {"", "null", 1000},
    }
    C := New[string, string]()
    c := C.(*cache[string, string])
    defer c.Close()

    for _, mock := range mocks {
        c.SetWithExpire(mock.key, mock.val, mock.ttl)
    }

    time.Sleep(wait)
    for _, mock := range mocks {
        val, found := c.Get(mock.key)
        if !found || val != mock.val {
            t.Fatalf("Unexpected value: %v (%v) to key: %v", val, found, mock.key)
        }
    }

    // update
    c.SetWithExpire(mocks[0].key, mocks[0].key+"foobar", 1000)
    val, found := c.Get(mocks[0].key)
    if !found || val == mocks[0].key {
        t.Fatalf("Unexpected value: %v (%v) to key: %v, want: %v", val, found, mocks[0].key, mocks[0].val)
    }

    // set struct
    type Foo struct {
        Name string
        Age  int
    }
    valStruct := Foo{
        "Bob",
        18,
    }
    valPtr := &valStruct
    c1 := New[string, Foo]()
    c2 := New[string, *Foo]()
    c1.SetWithExpire("foo", valStruct, 1000)
    c2.SetWithExpire("bar", valPtr, 1000)
    time.Sleep(wait)
    if val, found := c1.Get("foo"); !found || val != valStruct {
        t.Fatalf("Unexpected value: %v (%v) to key: %v, want: %v", val, found, "foo", valStruct)
    }
    if val, found := c2.Get("bar"); !found || val != valPtr {
        t.Fatalf("Unexpected value: %v (%v) to key: %v, want: %v", val, found, "foo", valPtr)
    }
}

// TestCacheMaxCap 测试Cache WithCapacity后的MaxCap是否符合预期
func TestCacheMaxCap(t *testing.T) {
    C := New(WithCapacity[string, string](2))
    c := C.(*cache[string, string])
    defer c.Close()

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
        c.SetWithExpire(mock.key, mock.val, mock.ttl)
    }

    time.Sleep(wait)
    if c.store.len() != 2 {
        t.Fatalf("unexpected Length:%d, want:%d", c.store.len(), 2)
    }
}

// TestCacheExpire 测试Cache的过期功能
func TestCacheExpire(t *testing.T) {
    C := New(WithCapacity[string, string](2), WithExpiration[string, string](1))
    c := C.(*cache[string, string])
    defer c.Close()

    c.Set("A", "a")
    c.Set("B", "b")

    time.Sleep(1)
    if c.store.len() != 0 {
        t.Fatalf("unexpected Length:%d, want:%d", c.store.len(), 0)
    }
    if val, found := c.Get("A"); found {
        t.Fatalf("unexpected expired value: %v to key: %v", val, "A")
    }
    if val, found := c.Get("B"); found {
        t.Fatalf("unexpected expired value: %v to key: %v", val, "B")
    }
}

// TestCacheSpecExpire 测试Cache的过期功能（对指定key）
func TestCacheSpecExpire(t *testing.T) {
    C := New[string, string](WithCapacity[string, string](2))
    c := C.(*cache[string, string])
    defer c.Close()

    mocks := []struct {
        key string
        val string
        ttl int64
    }{
        {"A", "a", 1},
        {"B", "b", 1},
        {"", "null", 1},
    }

    for _, mock := range mocks {
        c.SetWithExpire(mock.key, mock.val, mock.ttl)
    }

    time.Sleep(time.Second)
    if c.store.len() != 0 {
        t.Fatalf("unexpected Length:%d, want:%d", c.store.len(), 0)
    }
    for _, mock := range mocks {
        val, found := c.Get(mock.key)
        if found {
            t.Fatalf("unexpected expired value: %v to key: %v", val, mock.key)
        }
    }
}

// TestCacheGetWithLoad 测试Cache的WithLoad功能
func TestCacheGetWithLoad(t *testing.T) {
    m := map[string]string{
        "A": "a",
        "B": "b",
        "":  "null",
    }
    loadFunc := func(ctx context.Context, key string) (string, error) {
        if v, exist := m[key]; exist {
            return v, nil
        }
        return "", errors.New("key not exist")
    }

    C := New[string, string](WithLoad[string, string](loadFunc))
    c := C.(*cache[string, string])
    defer c.Close()

    type args struct {
        ctx context.Context
        key string
    }
    tests := []struct {
        name    string
        args    args
        want    interface{}
        wantErr bool
    }{
        {"A", args{context.Background(), "A"}, "a", false},
        {"B", args{context.Background(), "B"}, "b", false},
        {"", args{context.Background(), ""}, "null", false},
        {"unrecognized-key", args{context.Background(), "unregonizedKey"}, "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := c.GetWithLoad(tt.args.ctx, tt.args.key)
            time.Sleep(wait)
            got2, _ := c.Get(tt.args.key)
            if (err != nil) != tt.wantErr {
                t.Errorf("cache.GetWithLoad() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("cache.GetWithLoad() = %v, want %v", got, tt.want)
            }
            // Get获取
            if !reflect.DeepEqual(got2, tt.want) {
                t.Errorf("cache.Get() = %v, want %v", got2, tt.want)
            }
        })
    }
}

// TestCacheDel 测试Cache的删除功能
func TestCacheDel(t *testing.T) {
    C := New[string, string](WithCapacity[string, string](3))
    c := C.(*cache[string, string])
    defer c.Close()

    mocks := []struct {
        key string
        val string
        ttl int64
    }{
        {"A", "a", 10},
        {"B", "b", 10},
        {"", "null", 10},
    }

    for _, mock := range mocks {
        c.SetWithExpire(mock.key, mock.val, mock.ttl)
    }
    time.Sleep(wait)
    c.Del(mocks[0].key)
    time.Sleep(wait)

    if val, found := c.Get(mocks[0].key); found {
        t.Fatalf("unexpected deleted value: %v to key: %v", val, mocks[0].key)
    }

    for i := 1; i < len(mocks); i++ {
        if val, found := c.Get(mocks[i].key); !found || val != mocks[i].val {
            t.Fatalf("unexpected deleted value: %v (%v) to key: %v, want: %v", val, found, mocks[i].key, mocks[i].val)
        }
    }
}

// TestCacheClear 测试Cache的清空功能
func TestCacheClear(t *testing.T) {
    C := New[string, string]()
    c := C.(*cache[string, string])
    defer c.Close()
    for i := 0; i < 10; i++ {
        k := fmt.Sprint(i)
        v := fmt.Sprint(i)
        c.SetWithExpire(k, v, 10)
    }
    time.Sleep(wait)

    c.Clear()

    for i := 0; i < 10; i++ {
        k := fmt.Sprint(i)
        if _, found := c.Get(k); found {
            t.Fatalf("Shouldn't found value from clear cache")
        }
    }
    if c.store.len() != 0 {
        t.Fatalf("Length(%d) is not equal to 0, after clear", c.store.len())
    }
}

// BenchmarkCacheGet Benchmark Cache的Get功能
func BenchmarkCacheGet(b *testing.B) {
    k := "A"
    v := "a"

    c := New[string, string]()
    c.SetWithExpire(k, v, 100000)

    b.SetBytes(1)
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            c.Get(k)
        }
    })
}

// BenchmarkCacheGet Benchmark Cache的Set功能
func BenchmarkCacheSet(b *testing.B) {
    c := New[string, string]()
    rand.Seed(currentTime().Unix())

    b.SetBytes(1)
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            k := fmt.Sprint(rand.Int())
            v := k
            c.SetWithExpire(k, v, 10000)
        }
    })
}
