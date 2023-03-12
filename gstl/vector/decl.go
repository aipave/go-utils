package vector

// Iterator .
type Iterator[T any] interface {
	Next() bool
	Value() T
}

// Vector array container
type Vector[T any] interface {
	Begin() Iterator[T]

	RBegin() Iterator[T]

	Slice(begin, end int) []T

	Slices() []T

	// Get return value in pos
	Get(pos int) T

	Set(pos int, x T) bool

	Append(x ...T)

	Insert(pos int, x ...T)

	// Find return index of x in vector
	// -1 => not found
	// 0-N => index of x in the array
	Find(x T) int

	// FindBy custom defined filter
	// -1 => not found
	// 0-N => index ofx in the array
	FindBy(filter func(x T) bool) int

	Reverse()

	Size() int

	// Clear remove all elements
	Clear()

	// EraseAt remove index 'pos' from vector
	EraseAt(pos int) bool

	// Erase remove x from vector
	Erase(x T) bool

	// Sort by less
	Sort(less func(T, T) bool)

	// Front return first element of vector
	Front() T

	// Back return last element of vector
	Back() T

	// PopFront remove first element of vector
	PopFront() bool

	// PopBack remove last element of vector
	PopBack() bool

	Join(sep string) string
}
