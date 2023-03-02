package gcache

import (
	"container/list"
	"reflect"
	"sync"

	"github.com/alyu01/go-utils/gcast"

	"github.com/cespare/xxhash"
)

// store
// Is a concurrent safely used to store the key - value data storage. Temporarily use shard map implementation in this file
type store[K comparable, V any] interface {
	// get value of key
	get(K) (*list.Element, bool)
	// set key-value
	set(K, *list.Element)
	// delete key-value
	delete(K)
	// clear
	clear()
	// len Returns the size of the storage
	len() int
}

// newStore
// Returns the default implementation of the storage
func newStore[K comparable, V any]() store[K, V] {
	var val K
	switch reflect.ValueOf(val).Kind() {
	case reflect.Struct, reflect.Array, reflect.Chan, reflect.Interface, reflect.Ptr:
		return newLockedMap[K, V]()
	}

	return newShardedMap[K, V]()
}

const numShards uint64 = 256

// shardedMap Store the shard
type shardedMap[K comparable, V any] struct {
	shards []*lockedMap[K, V]
}

func newShardedMap[K comparable, V any]() *shardedMap[K, V] {
	sm := &shardedMap[K, V]{
		shards: make([]*lockedMap[K, V], int(numShards)),
	}
	for i := range sm.shards {
		sm.shards[i] = newLockedMap[K, V]()
	}
	return sm
}

func (sm *shardedMap[K, V]) get(key K) (*list.Element, bool) {

	return sm.shards[xxhash.Sum64String(gcast.ToString(key))&(numShards-1)].get(key)
}

func (sm *shardedMap[K, V]) set(key K, value *list.Element) {
	sm.shards[xxhash.Sum64String(gcast.ToString(key))&(numShards-1)].set(key, value)
}

func (sm *shardedMap[K, V]) delete(key K) {
	sm.shards[xxhash.Sum64String(gcast.ToString(key))&(numShards-1)].delete(key)
}

func (sm *shardedMap[K, V]) clear() {
	for i := uint64(0); i < numShards; i++ {
		sm.shards[i].clear()
	}
}

func (sm *shardedMap[K, V]) len() int {
	length := 0
	for i := uint64(0); i < numShards; i++ {
		length += sm.shards[i].len()
	}
	return length
}

// lockedMap Concurrent security Map
type lockedMap[K comparable, V any] struct {
	sync.RWMutex
	data map[K]*list.Element
}

func newLockedMap[K comparable, V any]() *lockedMap[K, V] {
	return &lockedMap[K, V]{
		data: make(map[K]*list.Element),
	}
}

func (m *lockedMap[K, V]) get(key K) (*list.Element, bool) {
	m.RLock()
	val, ok := m.data[key]
	m.RUnlock()
	return val, ok
}

func (m *lockedMap[K, V]) set(key K, value *list.Element) {
	m.Lock()
	m.data[key] = value
	m.Unlock()
}

func (m *lockedMap[K, V]) delete(key K) {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
}

func (m *lockedMap[K, V]) clear() {
	m.Lock()
	m.data = make(map[K]*list.Element)
	m.Unlock()
}

func (m *lockedMap[K, V]) len() int {
	return len(m.data)
}
