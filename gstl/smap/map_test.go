package smap

import (
	"fmt"
	"testing"
)

func TestSortedMap_Get(t *testing.T) {
	var m = NewMap[int, int]()

	fmt.Println("set [1, 1]", m.Set(1, 1))
	fmt.Println("set [1, 2]", m.Set(1, 2))

	fmt.Println("get 1", m.Get(1))

	if v, ok := m.GetOK(1); ok {
		fmt.Println("getOk 1", v)
	}

	if v, ok := m.GetOK(2); ok {
		fmt.Println("getOk 2", v)
	}

	var mm = NewMultiMap[int, int]()

	fmt.Println("set [1, 2]", mm.Set(1, 2))
	fmt.Println("set [1, 1]", mm.Set(1, 1))
	fmt.Println("set [1, 3]", mm.Set(1, 3))

	fmt.Println("get 1", mm.Get(1))

	for it := mm.GetFirst(1); it != nil; it = it.Next() {
		fmt.Println("get next 1", it.Value())
	}

	for it := mm.GetLast(1); it != nil; it = it.Next() {
		fmt.Println("get next 1", it.Value())
	}

	fmt.Println("get set", m.GetSet(1, 3))
	fmt.Println("get set", m.GetSet(1, 4))
}

func TestSortedMap_Keys(t *testing.T) {
	var m = NewMap[int, int]()

	for i := 1; i <= 10; i++ {
		m.Set(i, i*2)
	}

	fmt.Println("keys", m.Keys())
	fmt.Println("values", m.Values())
}

func TestSortedMap_Contains(t *testing.T) {
	var m = NewMap[int, int]()

	for i := 1; i <= 10; i++ {
		m.Set(i, i*2)
	}

	if m.Contains(3) {
		fmt.Println("duplicated key 3")
	}

	if m.Erase(3) {
		fmt.Println("removed key 3")
	}

	if !m.Contains(3) {
		fmt.Println("not contains 3")
	}
}

func TestSortedMap_Begin(t *testing.T) {
	var m = NewMap[int, int]()

	for i := 1; i <= 10; i++ {
		m.Set(i, i)
	}

	for it := m.Begin(); it != nil; it = it.Next() {
		fmt.Println("iter", it.Value())
	}

	for it := m.RBegin(); it != nil; it = it.Next() {
		fmt.Println("iter", it.Value())
	}

	for it := m.BeginFrom(3); it != nil; it = it.Next() {
		fmt.Println("iter", it.Value())
	}

	for it := m.RBeginFrom(3); it != nil; it = it.Next() {
		fmt.Println("iter", it.Value())
	}
}

func TestSortedMap_Intersect(t *testing.T) {
	var a = NewMap[int, int]()
	a.Set(1, 2)
	a.Set(2, 4)
	a.Set(4, 8)

	var b = NewMap[int, int]()
	b.Set(2, 4)
	b.Set(4, 8)
	b.Set(6, 12)

	var c = a.Intersect(b)
	c.Range(func(k, v int) bool {
		fmt.Println(k, v)
		return true
	})

	fmt.Println(a.Union(b).Keys())
}

func TestNewMap(t *testing.T) {
	var nm = NewMap[int, int]()
	fmt.Println(nm.Keys(), nm.Values())
}
