package vector

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// New create a vector[T] of type T
func New[T any](v []T) Vector[T] {
	var vec = &vector[T]{data: v}
	return vec
}

type vector[T any] struct {
	data []T
	mux  sync.RWMutex
}

// Sort by less
func (v *vector[T]) Sort(less func(a, b T) bool) {
	sort.Slice(v.data, func(i, j int) bool {
		return less(v.data[i], v.data[j])
	})
}

func (v *vector[T]) Size() int {
	return len(v.data)
}

func (v *vector[T]) Clear() {
	v.data = nil // clear
}

func (v *vector[T]) Append(value ...T) {
	v.mux.Lock()
	v.data = append(v.data, value...) // append not goroutine safe
	v.mux.Unlock()
}

// Slice 返回数组切片
// begin 0-n
// end -1 表示数组最后一个数
func (v *vector[T]) Slice(begin, end int) []T {
	if begin >= v.Size() {
		return nil
	}

	if end < 0 || end > v.Size() {
		end = v.Size()
	}

	return v.data[begin:end]
}

func (v *vector[T]) Slices() []T {
	return v.data
}

// Reverse 逆序函数
func (v *vector[T]) Reverse() {
	for i, j := 0, v.Size()-1; i < j; i, j = i+1, j-1 {
		v.data[i], v.data[j] = v.data[j], v.data[i]
	}
}

func (v *vector[T]) Get(pos int) (value T) {
	if pos < 0 || pos >= v.Size() {
		return
	}

	return v.data[pos]
}

// Set is a helper func to update val at 'pos' if it exists, otherwise nothing modified
func (v *vector[T]) Set(pos int, val T) bool {
	if pos >= 0 && pos < v.Size() {
		v.data[pos] = val
		return true
	}

	return false
}

func (v *vector[T]) Front() (val T) {
	if v.Size() == 0 {
		return
	}
	return v.data[0]
}

func (v *vector[T]) Back() (val T) {
	if v.Size() == 0 {
		return
	}

	index := v.Size() - 1
	return v.data[index]
}

func (v *vector[T]) PopBack() bool {
	if v.Size() > 0 {
		v.data = v.data[:v.Size()-1]
		return true
	}

	return false
}

func (v *vector[T]) PopFront() bool {
	if v.Size() > 0 {
		v.data = v.data[1:]
		return true
	}

	return false
}

func (v *vector[T]) EraseAt(pos int) bool {
	if pos < 0 || pos >= v.Size() {
		return false
	}

	for i := pos; i+1 < v.Size(); i++ {
		v.data[i] = v.data[i+1]
	}

	if pos < v.Size() {
		v.data = v.data[:v.Size()-1]
	}

	return true
}

// RBegin 逆序迭代
func (v *vector[T]) RBegin() Iterator[T] {
	return newReverseIterator(v, v.Size()-1)
}

// Begin 顺序迭代
func (v *vector[T]) Begin() Iterator[T] {
	return newIterator(v, 0)
}

// Erase is func for remote element from vector
// you need to get the iterator by Find before you call this method
func (v *vector[T]) Erase(x T) bool {
	return v.EraseAt(v.Find(x))
}

// Find parameter x must implements IEqual interface if its un-comparable type
func (v *vector[T]) Find(x T) int {
	for i := range v.data {
		if reflect.DeepEqual(v.data[i], x) {
			return i
		}
	}

	// find nothing
	return -1
}

func (v *vector[T]) FindBy(filter func(x T) bool) int {
	for i := range v.data {
		if filter(v.data[i]) {
			return i
		}
	}

	// find nothing
	return -1
}

// Join is a helper func to join vector value to string by 'sep'
func (v *vector[T]) Join(sep string) string {
	var elems []string
	for i := range v.data {
		elems = append(elems, fmt.Sprintf("%v", v.data[i]))
	}
	return strings.Join(elems, sep)
}

// Insert is helper func to batch insert elements to vector
func (v *vector[T]) Insert(pos int, val ...T) {
	if pos > v.Size() {
		pos = v.Size()
	}

	v.data = append(v.data, val...)

	for i := v.Size() - 1; i >= pos+len(val); i-- {
		v.data[i] = v.data[i-len(val)]
	}

	copy(v.data[pos:], val)
}

func (v *vector[T]) MarshalJSON() (b []byte, err error) {
	if v.data == nil {
		v.data = make([]T, 0)
	}

	return json.Marshal(v.data)
}

func (v *vector[T]) UnmarshalJSON(b []byte) (err error) {
	err = json.Unmarshal(b, &v.data)
	return
}
