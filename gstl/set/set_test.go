package set

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkSet_Insert(b *testing.B) {
	var s = NewSet[int]()

	rand.Seed(time.Now().UnixNano())
	for i := 1; i < b.N; i++ {
		s.Insert(rand.Int() % i)
		s.Erase(rand.Int() % i)
	}

	// fmt.Println(s.Size())
}

func TestSet_Begin(t *testing.T) {
	var s = NewSet[int]()

	for i := 1; i <= 100; i++ {
		s.Insert(i)
	}

	for it := s.Begin(); it != nil; it = it.Next() {
		fmt.Printf("%v ", it.Value())
	}

	fmt.Println()

	for it := s.RBegin(); it != nil; it = it.Next() {
		fmt.Printf("%v ", it.Value())
	}

	fmt.Println()
}

func TestSet_Contains(t *testing.T) {
	var s = NewSet[int]()

	for i := 1; i <= 100; i++ {
		s.Insert(i)
	}

	fmt.Println(s.Insert(98))
	fmt.Println(s.Erase(98))
	fmt.Println(s.Contains(98), s.Contains(99), s.Contains(101))
}

func TestSet_Insert(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	var s = NewExtSet(func(a, b User) int { return a.Age - b.Age })
	s.Insert(User{Name: "a", Age: 28})
	s.Insert(User{Name: "b", Age: 21})

	for it := s.Begin(); it != nil; it = it.Next() {
		fmt.Println(it.Value())
	}
}

func TestSet_Intersect(t *testing.T) {
	var a, b = NewSet[int](), NewSet[int]()

	for i := 1; i < 20; i++ {
		a.Insert(i)
	}

	for i := 16; i < 36; i++ {
		b.Insert(i)
	}

	a.Insert(100, 101, 102)

	fmt.Println("intersect", a.Intersect(b).Slices())
	fmt.Println("union", a.Union(b).Slices())
	fmt.Println("diff", a.Diff(b).Slices())
}

func TestSet_Slices(t *testing.T) {
	var a = NewSet[int64]()
	a.Insert(1029)
	a.Insert(1029)
	a.Insert(1028)
	a.Insert(1029)

	fmt.Println(a.Slices())
	for _, id := range a.Slices() {
		fmt.Println(id)
	}
}
