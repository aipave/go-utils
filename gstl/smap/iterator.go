package smap

import (
	"github.com/aipave/go-utils/gstl/compare"
	"github.com/aipave/go-utils/gstl/rbtree"
)

type iterator[K, V any] struct {
	current *rbtree.Node[K, V]
	sameKey *K
	keyCmp  compare.Less[K]
}

func newIterator[K, V any](node *rbtree.Node[K, V]) Iterator[K, V] {
	if node == nil {
		return nil
	}

	return iterator[K, V]{current: node}
}

func newSameIterator[K, V any](node *rbtree.Node[K, V], sameKey *K, keyCmp compare.Less[K]) Iterator[K, V] {
	if node == nil {
		return nil
	}

	return iterator[K, V]{current: node, sameKey: sameKey, keyCmp: keyCmp}
}

func (i iterator[K, V]) Next() Iterator[K, V] {
	if i.current == nil {
		return nil
	}

	next := i.current.Next()
	if next == nil {
		return nil
	}

	if i.sameKey != nil && i.keyCmp(next.Key(), *i.sameKey) != 0 {
		return nil
	}

	return iterator[K, V]{current: next, keyCmp: i.keyCmp}
}

func (i iterator[K, V]) Key() (k K) {
	if i.current != nil {
		return i.current.Key()
	}
	return
}

func (i iterator[K, V]) Value() (v V) {
	if i.current != nil {
		return i.current.Value()
	}
	return
}

type reverseIterator[K, V any] struct {
	current *rbtree.Node[K, V]
	sameKey *K
	keyCmp  compare.Less[K]
}

func newReverseIterator[K, V any](node *rbtree.Node[K, V]) Iterator[K, V] {
	if node == nil {
		return nil
	}

	return reverseIterator[K, V]{current: node}
}

func newReverseSameIterator[K, V any](node *rbtree.Node[K, V], sameKey *K, keyCmp compare.Less[K]) Iterator[K, V] {
	if node == nil {
		return nil
	}

	return reverseIterator[K, V]{current: node, sameKey: sameKey, keyCmp: keyCmp}
}

func (i reverseIterator[K, V]) Next() Iterator[K, V] {
	if i.current == nil {
		return nil
	}

	next := i.current.Next()
	if next == nil {
		return nil
	}

	if i.sameKey != nil && i.keyCmp(next.Key(), *i.sameKey) != 0 {
		return nil
	}

	return reverseIterator[K, V]{current: next, keyCmp: i.keyCmp}
}

func (i reverseIterator[K, V]) Key() (k K) {
	if i.current != nil {
		return i.current.Key()
	}
	return
}

func (i reverseIterator[K, V]) Value() (v V) {
	if i.current != nil {
		return i.current.Value()
	}
	return
}
