package set

import (
	"encoding/json"
	"fmt"

	"github.com/alyu01/go-utils/gstl/compare"
	"github.com/alyu01/go-utils/gstl/options"
	"github.com/alyu01/go-utils/gstl/rbtree"
)

type set[T any] struct {
    tree   *rbtree.RBTree[T, bool]
    mux    options.RWLock
    keyCmp compare.Less[T]
}

// NewSet 基本数据类型使用默认比较函数
func NewSet[T compare.Ordered](opts ...options.Option) Set[T] {
    var option = options.New()
    for _, fn := range opts {
        fn(&option)
    }

    return &set[T]{
        tree:   rbtree.New[T, bool](compare.OrderedLess[T]),
        mux:    option.Mux,
        keyCmp: compare.OrderedLess[T],
    }
}

// NewExtSet 非基本数据类型需要提供比较函数
func NewExtSet[T any](keyCmp compare.Less[T], opts ...options.Option) Set[T] {
    var option = options.New()
    for _, fn := range opts {
        fn(&option)
    }

    return &set[T]{
        tree:   rbtree.New[T, bool](keyCmp),
        mux:    option.Mux,
        keyCmp: keyCmp,
    }
}

func (s *set[T]) Size() int {
    return s.tree.Size()
}

func (s *set[T]) Clear() {
    s.mux.Lock()
    defer s.mux.Unlock()

    s.tree.Clear()
}

func (s *set[T]) Insert(list ...T) (allSuccess bool) {
    s.mux.Lock()
    defer s.mux.Unlock()

    for _, v := range list {
        node := s.tree.FindNode(v)
        if node == nil {
            s.tree.Insert(v, true)
        } else {
            allSuccess = false
        }
    }

    return
}

func (s *set[T]) Erase(v T) bool {
    s.mux.Lock()
    defer s.mux.Unlock()

    node := s.tree.FindNode(v)
    if node != nil {
        s.tree.Delete(node)
    }

    return node != nil
}

func (s *set[T]) BeginFrom(v T) Iterator[T] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindNode(v)
    return newIterator(node)
}

func (s *set[T]) Begin() Iterator[T] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    return newIterator(s.tree.Begin())
}

func (s *set[T]) RBegin() Iterator[T] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    return newReverseIterator(s.tree.RBegin())
}
func (s *set[T]) RBeginFrom(v T) Iterator[T] {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindNode(v)
    return newReverseIterator(node)
}

func (s *set[T]) Slices() (result []T) {
    for it := s.Begin(); it != nil; it = it.Next() {
        result = append(result, it.Value())
    }

    return
}

func (s *set[T]) Range(callback func(v T) bool) {
    for it := s.Begin(); it != nil; it = it.Next() {
        callback(it.Value())
    }
}

func (s *set[T]) Contains(v T) bool {
    s.mux.RLock()
    defer s.mux.RUnlock()

    node := s.tree.FindNode(v)
    return node != nil
}

// Intersect 交集
func (s *set[T]) Intersect(o Set[T]) Set[T] {
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

func (s *set[T]) Union(o Set[T]) Set[T] {
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

func (s *set[T]) Diff(o Set[T]) Set[T] {
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

func (s *set[T]) MarshalJSON() ([]byte, error) {
    s.mux.RLock()

    var list []T
    for it := s.Begin(); it != nil; it = it.Next() {
        list = append(list, it.Value())
    }

    s.mux.RUnlock()

    return json.Marshal(list)
}

func (s *set[T]) UnmarshalJSON(b []byte) (err error) {
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
