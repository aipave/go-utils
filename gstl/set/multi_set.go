package set

import (
	"encoding/json"
	"fmt"

	"github.com/alyu01/go-utils/gstl/compare"
	"github.com/alyu01/go-utils/gstl/options"
	"github.com/alyu01/go-utils/gstl/rbtree"
)

type multiSet[T any] struct {
    tree   *rbtree.RBTree[T, bool]
    mux    options.RWLock
    keyCmp compare.Less[T]
}

// NewMultiSet 基本数据类型使用默认比较函数
func NewMultiSet[T compare.Ordered](opts ...options.Option) MultiSet[T] {
    var option = options.New()
    for _, fn := range opts {
        fn(&option)
    }

    return &multiSet[T]{
        tree:   rbtree.New[T, bool](compare.OrderedLess[T]),
        mux:    option.Mux,
        keyCmp: compare.OrderedLess[T],
    }
}

// NewMultiExtSet 非基本数据类型需要提供比较函数
func NewMultiExtSet[T any](keyCmp compare.Less[T], opts ...options.Option) MultiSet[T] {
    var option = options.New()
    for _, fn := range opts {
        fn(&option)
    }

    return &multiSet[T]{
        tree:   rbtree.New[T, bool](keyCmp),
        mux:    option.Mux,
        keyCmp: keyCmp,
    }
}

func (s *multiSet[T]) Size() int {
    return s.tree.Size()
}

func (s *multiSet[T]) Clear() {
    s.mux.Lock()
    defer s.mux.Unlock()

    s.tree.Clear()
}

func (s *multiSet[T]) Range(callback func(v T) bool) {
    for it := s.Begin(); it != nil; it = it.Next() {
        callback(it.Value())
    }
}

func (s *multiSet[T]) Insert(list ...T) bool {
    s.mux.Lock()
    defer s.mux.Unlock()

    for _, v := range list {
        s.tree.Insert(v, true)
    }
    return true
}

func (s *multiSet[T]) Erase(v T) bool {
    s.mux.Lock()
    defer s.mux.Unlock()

    for {
        node := s.tree.FindNode(v)
        if node == nil {
            return true
        }
        s.tree.Delete(node)
    }
}

func (s *multiSet[T]) BeginFrom(v T) Iterator[T] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindFirstNode(v)
    return newIterator(node)
}

func (s *multiSet[T]) Begin() Iterator[T] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    return newIterator(s.tree.Begin())
}

func (s *multiSet[T]) RBegin() Iterator[T] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    return newReverseIterator(s.tree.RBegin())
}

func (s *multiSet[T]) RBeginFrom(v T) Iterator[T] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindFirstNode(v)
    return newReverseIterator(node)
}

func (s *multiSet[T]) Slices() (result []T) {
    for it := s.Begin(); it != nil; it = it.Next() {
        result = append(result, it.Value())
    }

    return
}

func (s *multiSet[T]) Contains(v T) bool {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindNode(v)
    return node != nil
}

func (s *multiSet[T]) MarshalJSON() ([]byte, error) {
    s.mux.RLock()

    var list []T
    for it := s.Begin(); it != nil; it = it.Next() {
        list = append(list, it.Value())
    }

    s.mux.RUnlock()

    return json.Marshal(list)
}

func (s *multiSet[T]) UnmarshalJSON(b []byte) (err error) {
    var list []T
    if !json.Valid(b) {
        return fmt.Errorf("input is not a valid json string")
    }

    err = json.Unmarshal(b, &list)
    if err != nil {
        return fmt.Errorf("unmarshal input to array struct err:%v", err)
    }

    s.Insert(list...)

    return
}

// Intersect 交集
func (s *multiSet[T]) Intersect(o Set[T]) Set[T] {
    var ans = NewExtSet(s.keyCmp)

    s.mux.RLock()

    for it := s.Begin(); it != nil; it = it.Next() {
        if o.Contains(it.Value()) {
            ans.Insert(it.Value())
        }
    }

    s.mux.RUnlock()

    return ans
}

func (s *multiSet[T]) Union(o Set[T]) Set[T] {
    var ans = NewExtSet(s.keyCmp)

    s.mux.RLock()

    for it := s.Begin(); it != nil; it = it.Next() {
        ans.Insert(it.Value())
    }

    s.mux.RUnlock()

    for it := o.Begin(); it != nil; it = it.Next() {
        ans.Insert(it.Value())
    }

    return ans
}

func (s *multiSet[T]) Diff(o Set[T]) Set[T] {
    var ans = NewExtSet(s.keyCmp)

    s.mux.RLock()

    for it := s.Begin(); it != nil; it = it.Next() {
        if !o.Contains(it.Value()) {
            ans.Insert(it.Value())
        }
    }

    s.mux.RUnlock()

    for it := o.Begin(); it != nil; it = it.Next() {
        if !s.Contains(it.Value()) {
            ans.Insert(it.Value())
        }
    }

    return ans
}
