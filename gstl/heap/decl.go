package heap

type Heap[T any] interface {
	Push(x ...T)
	Pop() T
	Size() int
}
