package heap

import (
    "github.com/alyu01/go-utils/gstl/compare"
)

type heap[T any] struct {
    arr    []T
    keyCmp compare.Less[T]
}

func New[T compare.Ordered]() Heap[T] {
    return &heap[T]{keyCmp: compare.OrderedLess[T]}
}

func NewExt[T any](cmp compare.Less[T]) Heap[T] {
    return &heap[T]{keyCmp: cmp}
}

func (h *heap[T]) Push(x ...T) {
    h.arr = append(h.arr, x...)

    for i := h.Size() - len(x) - 1; i < h.Size(); i++ {
        h.up(i)
    }
}

func (h *heap[T]) Pop() (v T) {
    n := h.Size() - 1
    v = h.arr[0]
    h.swap(0, n)
    h.down(0, n)
    h.arr = h.arr[:n]
    return
}

func (h *heap[T]) Remove(i int) any {
    n := h.Size() - 1
    if n != i {
        h.swap(i, n)
        if !h.down(i, n) {
            h.up(i)
        }
    }
    return h.Pop()
}

func (h *heap[T]) Fix(i int) {
    if !h.down(i, h.Size()) {
        h.up(i)
    }
}

func (h *heap[T]) up(j int) {
    for {
        i := (j - 1) / 2 // parent
        if i == j || h.keyCmp(h.arr[j], h.arr[i]) >= 0 {
            break
        }
        h.swap(i, j)
        j = i
    }
}

func (h *heap[T]) down(i0, n int) bool {
    i := i0
    for {
        j1 := 2*i + 1
        if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
            break
        }
        j := j1 // left child
        if j2 := j1 + 1; j2 < n && h.keyCmp(h.arr[j2], h.arr[j1]) < 0 {
            j = j2 // = 2*i + 2  // right child
        }
        if h.keyCmp(h.arr[j], h.arr[i]) >= 0 {
            break
        }
        h.swap(i, j)
        i = j
    }
    return i > i0
}

func (h *heap[T]) swap(i, j int) {
    h.arr[i], h.arr[j] = h.arr[j], h.arr[i]
}

func (h *heap[T]) Size() int {
    return len(h.arr)
}
