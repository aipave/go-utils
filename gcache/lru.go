package gcache

import (
	"container/list"
	"sync"
	"time"
)

// entry Storage entity
type entry[K comparable, V any] struct {
	mux        sync.RWMutex
	key        K
	value      V
	expireTime time.Time
}

func getEntry[K comparable, V any](ele *list.Element) *entry[K, V] {
	return ele.Value.(*entry[K, V])
}

func setEntry[K comparable, V any](ele *list.Element, ent *entry[K, V]) {
	ele.Value = ent
}

// lru The concurrent security lru queues
type lru[K comparable, V any] struct {
	ll       *list.List
	store    store[K, V]
	capacity int
}

func newLRU[K comparable, V any](capacity int, store store[K, V]) *lru[K, V] {
	return &lru[K, V]{
		ll:       list.New(),
		store:    store,
		capacity: capacity,
	}
}

func (l *lru[K, V]) add(ent *entry[K, V]) *entry[K, V] {
	val, ok := l.store.get(ent.key)
	if ok {
		setEntry(val, ent)
		l.ll.MoveToFront(val)
		return nil
	}
	if l.capacity <= 0 || l.ll.Len() < l.capacity {
		ele := l.ll.PushFront(ent)
		l.store.set(ent.key, ele)
		return nil
	}
	// After the expiration of the lru, replace the last element
	ele := l.ll.Back()
	if ele == nil {
		return ent
	}
	victimEnt := getEntry[K, V](ele)
	setEntry(ele, ent)
	l.ll.MoveToFront(ele)

	l.store.delete(victimEnt.key)
	l.store.set(ent.key, ele)
	return victimEnt
}

func (l *lru[K, V]) hit(ele *list.Element) {
	l.ll.MoveToFront(ele)
}

func (l *lru[K, V]) push(elements []*list.Element) {
	for _, ele := range elements {
		l.ll.MoveToFront(ele)
	}
}

func (l *lru[K, V]) delete(key K) {
	value, ok := l.store.get(key)
	if !ok {
		return
	}
	l.ll.Remove(value)
	l.store.delete(key)
}

func (l *lru[K, V]) len() int {
	return l.ll.Len()
}

func (l *lru[K, V]) clear() {
	l.ll = list.New()
}
