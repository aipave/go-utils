package gcache

import (
	"container/list"
)

// Storage policy
type policy[K comparable, V any] interface {
	// add  an elem.
	add(ent *entry[K, V]) *entry[K, V]
	// hit processes a hit on an element
	hit(elements *list.Element)
	// push processes a batch of accessed elements
	push(elements []*list.Element)
	// delete
	delete(key K)
	// clear
	clear()
}

func newPolicy[K comparable, V any](capacity int, store store[K, V]) policy[K, V] {
	return newLRU[K, V](capacity, store)
}
