package vector

func newIterator[T any](vec *vector[T], pos int) Iterator[T] {
	return &iterator[T]{vec: vec, current: vec.Get(pos), pos: pos}
}

type iterator[T any] struct {
	vec     *vector[T]
	current T
	pos     int
}

func (v *iterator[T]) Next() bool {
	if v.pos >= v.vec.Size() {
		return false
	}

	v.current = v.vec.Get(v.pos)
	v.pos++
	return true
}

func (v *iterator[T]) Value() T {
	return v.current
}

func newReverseIterator[T any](vec *vector[T], pos int) Iterator[T] {
	return &reverseIterator[T]{vec: vec, current: vec.Get(pos), pos: pos}
}

type reverseIterator[T any] struct {
	vec     *vector[T]
	current T
	pos     int
}

func (v *reverseIterator[T]) Next() bool {
	if v.pos < 0 {
		return false
	}

	v.current = v.vec.Get(v.pos)
	v.pos--
	return true
}

func (v reverseIterator[T]) Value() T {
	return v.current
}
