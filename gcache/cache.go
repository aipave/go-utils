package gcache

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

const (
	// The maximum capacity of the cache
	maxCapacity = 1 << 30

	ringBufSize = 16

	// Asynchronous metadata in cache operations
	setBufSize    = 1 << 15
	updateBufSize = 1 << 15
	delBufSize    = 1 << 15
	expireBufSize = 1 << 15

	// The default expiration time 60 s
	ttl = 60 * time.Second
)

// currentTime for convenience when testing
var currentTime = time.Now

// Cache is a local support expiration time K - V memory storage
type Cache[K comparable, V any] interface {
	// Get value of key, if key not exist or expired, return false
	Get(key K) (V, bool)
	// GetWithLoad  Get value of key, If the key doesn't exist, using user-defined filling load data returns, and caching
	GetWithLoad(ctx context.Context, key K) (V, error)
	// Set key, value
	Set(key K, value V) bool
	// SetWithExpire set key，value, ttl(s)
	SetWithExpire(key K, value V, ttl int64) bool
	// Del key
	Del(key K)
	// Clear all queue and cache
	Clear()
	// Close cache
	Close()
}

// cache K-V
type cache[K comparable, V any] struct {
	capacity int
	store    store[K, V]

	// Access and exit strategies of key
	policy policy[K, V]

	getBuf     *ringBuffer
	elementsCh chan []*list.Element

	setBuf    chan *entry[K, V]
	updateBuf chan *list.Element
	delBuf    chan K
	expireBuf chan K

	g    groupType[K, V]
	load LoadFunc[K, V]

	ttl time.Duration
	// delete expire key's task-queue
	expireQueue *expireQueue[K, V]

	stop chan struct{}
}

// Load key's value for populate the cache
type LoadFunc[K comparable, V any] func(ctx context.Context, key K) (V, error)

// call (With reference to singleflight) For user-defined data according to the load
type callType[V any] struct {
	wg  sync.WaitGroup
	val V
	err error
}

// group (With reference to singleflight) For user-defined data according to the load
type groupType[K comparable, V any] struct {
	mu sync.Mutex         // protects m
	m  map[K]*callType[V] // lazily initialized
}

// Option para tools func
type Option[K comparable, V any] func(*cache[K, V])

// WithCapacity set key's max cpa.
func WithCapacity[K comparable, V any](capacity int) Option[K, V] {
	if capacity > maxCapacity {
		capacity = maxCapacity
	}
	if capacity <= 0 {
		capacity = 1
	}
	return func(c *cache[K, V]) {
		c.capacity = capacity
	}
}

// WithExpiration set ttl(s), default time is 60s.
func WithExpiration[K comparable, V any](ttl time.Duration) Option[K, V] {
	if ttl <= 0 {
		ttl = 1
	}
	return func(c *cache[K, V]) {
		c.ttl = ttl
	}
}

// WithLoad Set the custom data loading function
func WithLoad[K comparable, V any](f LoadFunc[K, V]) Option[K, V] {
	return func(c *cache[K, V]) {
		c.load = f
	}
}

// New cache obj.
func New[K comparable, V any](opts ...Option[K, V]) Cache[K, V] {
	// initial cache with default value
	cache := &cache[K, V]{
		capacity: maxCapacity,
		store:    newStore[K, V](),

		elementsCh: make(chan []*list.Element, 3),
		setBuf:     make(chan *entry[K, V], setBufSize),
		updateBuf:  make(chan *list.Element, updateBufSize),
		delBuf:     make(chan K, delBufSize),
		expireBuf:  make(chan K, expireBufSize),

		ttl:         ttl,
		expireQueue: newExpireQueue[K, V](time.Second, 60),
		stop:        make(chan struct{}),
	}
	// set cache with custom para.
	for _, opt := range opts {
		opt(cache)
	}

	cache.policy = newPolicy[K, V](cache.capacity, cache.store)
	cache.getBuf = newRingBuffer(cache, ringBufSize)

	go cache.processEntries()

	return cache
}

// processEntries Asynchronous processing to the operation of the cache
func (c *cache[K, V]) processEntries() {
	for {
		select {
		case elements := <-c.elementsCh:
			c.access(elements)
		case ent := <-c.setBuf:
			c.add(ent)
		case ele := <-c.updateBuf:
			c.update(ele)
		case key := <-c.delBuf:
			c.delete(key)
		case key := <-c.expireBuf:
			c.expire(key)
		case <-c.stop:
			return
		}
	}
}

// Get value of key, if key not exist or expired, return false
func (c *cache[K, V]) Get(key K) (result V, ok bool) {
	if c == nil {
		return
	}

	value, hit := c.store.get(key)
	if hit {
		ent := getEntry[K, V](value)
		ent.mux.RLock()
		defer ent.mux.RUnlock()
		if ent.expireTime.Before(currentTime()) {
			return
		}
		c.getBuf.push(value)
		return ent.value, true
	}

	return
}

// GetWithLoad  Get value of key, If the key doesn't exist, using user-defined filling load data returns, and caching
func (c *cache[K, V]) GetWithLoad(ctx context.Context, key K) (result V, err error) {
	if c.load == nil {
		return result, errors.New("undefined LoadFunc in cache")
	}

	value, found := c.Get(key)
	if found {
		return value, nil
	}

	return c.loadData(ctx, key, c.load)
}

// Set key, value
func (c *cache[K, V]) Set(key K, value V) bool {
	return c.SetWithExpire(key, value, int64(c.ttl.Seconds()))
}

// SetWithExpire set key，value，time to live (s).
func (c *cache[K, V]) SetWithExpire(key K, value V, ttl int64) bool {
	if c == nil {
		return false
	}
	expireTime := currentTime().Add(time.Duration(ttl) * time.Second)
	ent := &entry[K, V]{
		key:        key,
		value:      value,
		expireTime: expireTime,
	}

	val, hit := c.store.get(key)
	if hit {
		// If the key exists, immediately update stored in the latest value, prevent the Get access to the dirty data
		oldEnt := getEntry[K, V](val)

		oldEnt.mux.Lock()
		oldEnt.value = value
		oldEnt.expireTime = expireTime
		oldEnt.mux.Unlock()

		select {
		case c.updateBuf <- val:
		default:
		}
		c.expireQueue.update(ent.key, ent.expireTime, c.afterExpire(ent.key))

		return true
	}
	// The new key and the value. Extreme cases, Set operation is no guarantee of success, because the asynchronous processing, when the load is bigger, the Set result ms level possible delays
	select {
	case c.setBuf <- ent:
		return true
	default:
		return false
	}
}

// Del key
func (c *cache[K, V]) Del(key K) {
	if c == nil {
		return
	}
	c.delBuf <- key
}

// Clean all queues and caching, atomic operation, should be called after didn't Get and Set operation
func (c *cache[K, V]) Clear() {
	// Block until processEntries goroutine to an end
	c.stop <- struct{}{}

	c.elementsCh = make(chan []*list.Element, 3)
	c.setBuf = make(chan *entry[K, V], setBufSize)
	c.updateBuf = make(chan *list.Element, updateBufSize)
	c.delBuf = make(chan K, delBufSize)
	c.expireBuf = make(chan K, expireBufSize)

	c.store.clear()
	c.policy.clear()
	c.expireQueue.clear()

	// restart processEntries goroutine
	go c.processEntries()
}

// Close cache
func (c *cache[K, V]) Close() {
	// block until processEntries goroutine's ending
	c.stop <- struct{}{}
	close(c.stop)

	close(c.elementsCh)
	close(c.setBuf)
	close(c.updateBuf)
	close(c.delBuf)
	close(c.expireBuf)

	c.expireQueue.stop()
}

// access Is an asynchronous invocation, handle read operations
func (c *cache[K, V]) access(elements []*list.Element) {
	c.policy.push(elements)
}

// add Is an asynchronous invocation, handle write operations
func (c *cache[K, V]) add(ent *entry[K, V]) {
	// Store the new key - value, after reaching maximum size, return to entry to be eliminated
	key := ent.key
	expireTime := ent.expireTime
	victimEnt := c.policy.add(ent)
	c.expireQueue.add(key, expireTime, c.afterExpire(key))

	// Deleted from the queue date entry to be eliminated
	if victimEnt != nil {
		c.expireQueue.remove(victimEnt.key)
	}
}

// update Is an asynchronous invocation, processing the update operation
func (c *cache[K, V]) update(ele *list.Element) {
	c.policy.hit(ele)
}

// delete Is an asynchronous call, delete operation
func (c *cache[K, V]) delete(key K) {
	c.policy.delete(key)
	c.expireQueue.remove(key)
}

// expire Be an asynchronous invocation, outdated data processing
func (c *cache[K, V]) expire(key K) {
	c.policy.delete(key)
}

// loadData Use the user-defined functions to load and cached data
func (c *cache[K, V]) loadData(ctx context.Context, key K, load func(context.Context, K) (V, error)) (V, error) {
	// Reference singleflight implementation, the cache breakdown
	c.g.mu.Lock()
	if c.g.m == nil {
		c.g.m = make(map[K]*callType[V])
	}
	if call, ok := c.g.m[key]; ok {
		c.g.mu.Unlock()
		call.wg.Wait()
		return call.val, call.err
	}
	call := new(callType[V])
	call.wg.Add(1)
	c.g.m[key] = call
	c.g.mu.Unlock()

	call.val, call.err = load(ctx, key)
	if call.err == nil {
		c.Set(key, call.val)
	}

	c.g.mu.Lock()
	call.wg.Done()
	delete(c.g.m, key)
	c.g.mu.Unlock()

	return call.val, call.err
}

// afterExpire Return the key after the expiration of the callback
func (c *cache[K, V]) afterExpire(key K) func() {
	return func() {
		select {
		case c.expireBuf <- key:
		default:
		}
	}
}

// push Will read requests batch into the channel
func (c *cache[K, V]) push(elements []*list.Element) bool {
	if len(elements) == 0 {
		return true
	}
	select {
	case c.elementsCh <- elements:
		return true
	default:
		return false
	}
}
