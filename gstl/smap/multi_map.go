package smap

import (
    "github.com/alyu01/go-utils/gstl/compare"
    "github.com/alyu01/go-utils/gstl/options"
    "github.com/alyu01/go-utils/gstl/rbtree"
)

type sortedMultiMap[K, V any] struct {
    tree   *rbtree.RBTree[K, V]
    mux    options.RWLock
    keyCmp compare.Less[K]
}

// NewMultiMap
func NewMultiMap[K compare.Ordered, V any](opts ...options.Option) MultiMap[K, V] {
    var option = options.New()
    for _, fn := range opts {
        fn(&option)
    }

    return &sortedMultiMap[K, V]{
        tree:   rbtree.New[K, V](compare.OrderedLess[K]),
        mux:    option.Mux,
        keyCmp: compare.OrderedLess[K],
    }
}

// NewMultiExtMap
func NewMultiExtMap[K, V any](keyCmp compare.Less[K], opts ...options.Option) MultiMap[K, V] {
    var option = options.New()
    for _, fn := range opts {
        fn(&option)
    }

    return &sortedMultiMap[K, V]{
        tree:   rbtree.New[K, V](keyCmp),
        mux:    option.Mux,
        keyCmp: keyCmp,
    }
}

func (s *sortedMultiMap[K, V]) Size() int {
    return s.tree.Size()
}

func (s *sortedMultiMap[K, V]) Clear() {
    s.mux.Lock()
    defer s.mux.Unlock()

    s.tree.Clear()
}

func (s *sortedMultiMap[K, V]) Get(key K) (result V) {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindNode(key)
    if node != nil {
        return node.Value()
    }

    return
}

func (s *sortedMultiMap[K, V]) GetOK(key K) (result V, b bool) {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindNode(key)
    if node != nil {
        return node.Value(), true
    }

    return
}

func (s *sortedMultiMap[K, V]) GetFirst(key K) Iterator[K, V] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindFirstNode(key)
    if node != nil {
        return newSameIterator(node, &key, s.keyCmp)
    }

    return nil
}

func (s *sortedMultiMap[K, V]) GetLast(key K) Iterator[K, V] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindLastNode(key)
    if node != nil {
        return newReverseSameIterator(node, &key, s.keyCmp)
    }

    return nil
}

func (s *sortedMultiMap[K, V]) Set(key K, value V) bool {
    s.mux.Lock()
    defer s.mux.Unlock()

    node := s.tree.FindNode(key)
    s.tree.Insert(key, value)

    return node == nil
}

func (s *sortedMultiMap[K, V]) Erase(v K) (ok bool) {
    s.mux.Lock()
    defer s.mux.Unlock()

    for {
        node := s.tree.FindNode(v)
        if node == nil {
            break
        }
        s.tree.Delete(node)
        ok = true
    }

    return
}

func (s *sortedMultiMap[K, V]) First() (v V) {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.Begin()
    if node != nil {
        return node.Value()
    }

    return
}

func (s *sortedMultiMap[K, V]) Last() (v V) {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.RBegin()
    if node != nil {
        return node.Value()
    }

    return
}

func (s *sortedMultiMap[K, V]) BeginFrom(v K) Iterator[K, V] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindFirstNode(v)
    return newIterator(node)
}

func (s *sortedMultiMap[K, V]) Begin() Iterator[K, V] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    return newIterator(s.tree.Begin())
}

func (s *sortedMultiMap[K, V]) RBegin() Iterator[K, V] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    return newReverseIterator(s.tree.RBegin())
}
func (s *sortedMultiMap[K, V]) RBeginFrom(v K) Iterator[K, V] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindLastNode(v)
    return newReverseIterator(node)
}

func (s *sortedMultiMap[K, V]) Keys() (result []K) {
    for it := s.Begin(); it != nil; it = it.Next() {
        result = append(result, it.Key())
    }

    return
}

func (s *sortedMultiMap[K, V]) Values() (result []V) {
    for it := s.Begin(); it != nil; it = it.Next() {
        result = append(result, it.Value())
    }

    return
}

func (s *sortedMultiMap[K, V]) Contains(v K) bool {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindNode(v)
    return node != nil
}
