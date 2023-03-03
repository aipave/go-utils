package set

type Iterator[T any] interface {
	Next() Iterator[T]
	Value() T
}

type Set[T any] interface {
	Size() int

	Slices() []T

	// Clear remove all elements from set
	Clear()

	// Insert add elements to set, return false if any node exist
	Insert(v ...T) bool

	// Erase delete element from set, return true if node exist
	Erase(v T) bool

	Contains(v T) bool

	Range(func(T) bool)

	BeginFrom(v T) Iterator[T]

	Begin() Iterator[T]

	RBeginFrom(v T) Iterator[T]

	RBegin() Iterator[T]

	Intersect(Set[T]) Set[T]

	Union(Set[T]) Set[T]

	Diff(Set[T]) Set[T]
}

type MultiSet[T any] interface {
	Set[T]
}
