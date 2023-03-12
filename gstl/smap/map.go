package smap

import (
	"github.com/aipave/go-utils/gstl/compare"
	"github.com/aipave/go-utils/gstl/options"
	"github.com/aipave/go-utils/gstl/rbtree"
)

type sortedMap[K, V any] struct {
	tree   *rbtree.RBTree[K, V]
	mux    options.RWLock
	keyCmp compare.Less[K]
}

// NewMap Primitive data types use default comparison functions
func NewMap[K compare.Ordered, V any](opts ...options.Option) Map[K, V] {
	var option = options.New()
	for _, fn := range opts {
		fn(&option)
	}

	return &sortedMap[K, V]{
		tree:   rbtree.New[K, V](compare.OrderedLess[K]),
		mux:    option.Mux,
		keyCmp: compare.OrderedLess[K],
	}
}

// NewExtMap Non-basic data types need to provide comparison functions
func NewExtMap[K, V any](keyCmp compare.Less[K], opts ...options.Option) Map[K, V] {
	var option = options.New()
	for _, fn := range opts {
		fn(&option)
	}

	return &sortedMap[K, V]{
		tree:   rbtree.New[K, V](keyCmp),
		mux:    option.Mux,
		keyCmp: keyCmp,
	}
}

func (s *sortedMap[K, V]) Size() int {
	return s.tree.Size()
}

func (s *sortedMap[K, V]) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.tree.Clear()
}

func (s *sortedMap[K, V]) Get(key K) (result V) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	node := s.tree.FindNode(key)
	if node != nil {
		return node.Value()
	}

	return
}

func (s *sortedMap[K, V]) GetOK(key K) (result V, b bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	node := s.tree.FindNode(key)
	if node != nil {
		return node.Value(), true
	}

	return
}

func (s *sortedMap[K, V]) Set(key K, value V) (allSuccess bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	node := s.tree.FindNode(key)
	if node != nil {
		node.SetValue(value)
	} else {
		s.tree.Insert(key, value)
	}

	return node == nil
}

func (s *sortedMap[K, V]) GetSet(key K, value V) (result V) {
	s.mux.Lock()
	defer s.mux.Unlock()

	node := s.tree.FindNode(key)
	if node != nil {
		result = node.Value()
		node.SetValue(value)
	} else {
		s.tree.Insert(key, value)
	}

	return
}

func (s *sortedMap[K, V]) GetSetOK(key K, value V) (result V, b bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	node := s.tree.FindNode(key)
	if node != nil {
		result, b = node.Value(), true
		node.SetValue(value)
	} else {
		s.tree.Insert(key, value)
	}

	return
}

func (s *sortedMap[K, V]) Erase(v K) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	node := s.tree.FindNode(v)
	if node != nil {
		s.tree.Delete(node)
	}

	return node != nil
}

func (s *sortedMap[K, V]) First() (v V) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	node := s.tree.Begin()
	if node != nil {
		return node.Value()
	}

	return
}

func (s *sortedMap[K, V]) Last() (v V) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	node := s.tree.RBegin()
	if node != nil {
		return node.Value()
	}

	return
}

func (s *sortedMap[K, V]) BeginFrom(v K) Iterator[K, V] {
	s.mux.RLock()
	defer s.mux.RUnlock()

	node := s.tree.FindNode(v)
	return newIterator(node)
}

func (s *sortedMap[K, V]) Begin() Iterator[K, V] {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return newIterator(s.tree.Begin())
}

func (s *sortedMap[K, V]) RBegin() Iterator[K, V] {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return newReverseIterator(s.tree.RBegin())
}
func (s *sortedMap[K, V]) RBeginFrom(v K) Iterator[K, V] {
	s.mux.RLock()
	defer s.mux.RUnlock()

	node := s.tree.FindNode(v)
	return newReverseIterator(node)
}

func (s *sortedMap[K, V]) Keys() (result []K) {
	for it := s.Begin(); it != nil; it = it.Next() {
		result = append(result, it.Key())
	}

	return
}

func (s *sortedMap[K, V]) Values() (result []V) {
	for it := s.Begin(); it != nil; it = it.Next() {
		result = append(result, it.Value())
	}

	return
}

func (s *sortedMap[K, V]) Contains(v K) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	node := s.tree.FindNode(v)
	return node != nil
}

// Intersect 交集
func (s *sortedMap[K, V]) Intersect(o Map[K, V]) Map[K, V] {
	var ans = NewExtMap[K, V](s.keyCmp)

	s.mux.RLock()

	for it := s.Begin(); it != nil; it = it.Next() {
		if o.Contains(it.Key()) {
			ans.Set(it.Key(), it.Value())
		}
	}

	s.mux.RUnlock()

	return ans
}

func (s *sortedMap[K, V]) Union(o Map[K, V]) Map[K, V] {
	var ans = NewExtMap[K, V](s.keyCmp)

	s.mux.RLock()

	for it := s.Begin(); it != nil; it = it.Next() {
		ans.Set(it.Key(), it.Value())
	}

	s.mux.RUnlock()

	for it := o.Begin(); it != nil; it = it.Next() {
		ans.Set(it.Key(), it.Value())
	}

	return ans
}

func (s *sortedMap[K, V]) Diff(o Map[K, V]) Map[K, V] {
	var ans = NewExtMap[K, V](s.keyCmp)

	s.mux.RLock()

	for it := s.Begin(); it != nil; it = it.Next() {
		if !o.Contains(it.Key()) {
			ans.Set(it.Key(), it.Value())
		}
	}

	s.mux.RUnlock()

	for it := o.Begin(); it != nil; it = it.Next() {
		if !s.Contains(it.Key()) {
			ans.Set(it.Key(), it.Value())
		}
	}

	return ans
}

func (s *sortedMap[K, V]) Range(callback func(k K, v V) bool) {
	for it := s.Begin(); it != nil; it = it.Next() {
		callback(it.Key(), it.Value())
	}
}
