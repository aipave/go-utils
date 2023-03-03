package smap

type Iterator[K, V any] interface {
	Next() Iterator[K, V]
	Key() K
	Value() V
}

type Map[K, V any] interface {
	// Get return default value of type V if key not exist, otherwise return the value of key
	Get(key K) V

	// GetOK return true if key exist, otherwise return false
	GetOK(key K) (V, bool)

	// Set return true if key not exist, otherwise return false
	Set(key K, value V) bool

	// GetSet return the old value of key and set it to the new value
	GetSet(key K, value V) V

	// GetSetOK same as GetSet, but return false if old key does not exist
	GetSetOK(key K, value V) (V, bool)

	Contains(key K) bool

	// Erase return true if key exist, otherwise return false
	Erase(key K) bool

	First() V

	Last() V

	Begin() Iterator[K, V]

	BeginFrom(key K) Iterator[K, V]

	RBegin() Iterator[K, V]

	RBeginFrom(key K) Iterator[K, V]

	Keys() []K

	Values() []V

	Clear()

	Size() int

	Range(func(K, V) bool)

	Intersect(p Map[K, V]) Map[K, V]

	Union(Map[K, V]) Map[K, V]

	Diff(Map[K, V]) Map[K, V]
}

type MultiMap[K, V any] interface {
	// Get return the node which first visited in the tree
	Get(key K) V

	// GetOK return true if key exist, otherwise return false
	GetOK(key K) (V, bool)

	// GetFirst return the left node of same key, iterator next will range from left to right
	GetFirst(key K) Iterator[K, V]

	// GetLast return the right node of same key, iterator next will range from right to left
	GetLast(key K) Iterator[K, V]

	// Set return true if key not exist, otherwise return false
	Set(key K, value V) bool

	Contains(key K) bool

	// Erase return true if key exist, otherwise return false
	Erase(key K) bool

	First() V

	Last() V

	Begin() Iterator[K, V]

	BeginFrom(key K) Iterator[K, V]

	RBegin() Iterator[K, V]

	RBeginFrom(key K) Iterator[K, V]

	Keys() []K

	Values() []V

	Clear()

	Size() int
}
