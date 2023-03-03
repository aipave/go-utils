package set

import (
    "github.com/alyu01/go-utils/gstl/rbtree"
)

type iterator[T any] struct {
    node *rbtree.Node[T, bool]
}

func newIterator[T any](node *rbtree.Node[T, bool]) Iterator[T] {
    if node == nil {
        return nil
    }

    return iterator[T]{node: node}
}

func (iter iterator[T]) Next() Iterator[T] {
    if iter.node == nil {
        return nil
    }

    nextNode := iter.node.Next()
    if nextNode == nil {
        return nil
    }

    return iterator[T]{node: nextNode}
}

func (iter iterator[T]) Value() (v T) {
    if iter.node != nil {
        return iter.node.Key()
    }

    return
}

type reverseIterator[T any] struct {
    node *rbtree.Node[T, bool]
}

func newReverseIterator[T any](node *rbtree.Node[T, bool]) Iterator[T] {
    if node == nil {
        return nil
    }

    return reverseIterator[T]{node: node}
}

func (iter reverseIterator[T]) Next() Iterator[T] {
    if iter.node == nil {
        return nil
    }

    nextNode := iter.node.Next()
    if nextNode == nil {
        return nil
    }

    return reverseIterator[T]{node: nextNode}
}

func (iter reverseIterator[T]) Value() (v T) {
    if iter.node != nil {
        return iter.node.Key()
    }

    return
}
