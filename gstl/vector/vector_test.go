package vector

import (
	"encoding/json"
	"fmt"
	"testing"
)

type UnComparable struct {
	Key   int
	Value []int
}

func (UnComparable) Equal(a, b UnComparable) bool {
	return a.Key == b.Key
}

func TestVector_At(t *testing.T) {
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)
	fmt.Println(v.Get(8), v.Get(100))
}

func TestVector_Back(t *testing.T) {
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)
	fmt.Println(v.Back())
}

func TestVector_EraseAt(t *testing.T) {
	// comparable type
	fmt.Println("=========== comparable ==============")
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)
	v.EraseAt(3)
	fmt.Println(v)
}

func TestVector_Front(t *testing.T) {
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)

	fmt.Println(v.Front())
}

func TestVector_PopBack(t *testing.T) {
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)
	fmt.Println(v.PopBack())
}

func TestVector_PushBack(t *testing.T) {
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)

	v.Append(10, 11, 12)
	fmt.Println(v.Slices())
}

func TestVector_Reverse(t *testing.T) {
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)

	v.Reverse()
	fmt.Println(v.Slices())
}

func TestVector_Size(t *testing.T) {
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)

	fmt.Println(v.Size())
}

func TestVector_Slice(t *testing.T) {
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)

	fmt.Println(v.Slice(3, 100))
}

func TestVector_Sort(t *testing.T) {
	// comparable type
	fmt.Println("=========== comparable ==============")
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)

	fmt.Println("before", v.Slice(0, v.Size()))
	v.Sort(func(a, b int) bool {
		return a > b
	})
	fmt.Println("after", v.Slice(0, v.Size()))

	// un-comparable type
	fmt.Println("\n=========== un-comparable ==============")
	var vv = New([]UnComparable{{Key: 1, Value: []int{11}}, {Key: 3, Value: []int{33}}, {Key: 2, Value: []int{22}}, {Key: 4, Value: []int{44}}})
	fmt.Println("before", vv.Slice(0, vv.Size()))
	vv.Sort(func(a, b UnComparable) bool {
		return a.Key < b.Key
	})
	fmt.Println("after", vv.Slice(0, vv.Size()))
}

func TestVector_Find(t *testing.T) {
	// comparable type
	fmt.Println("=========== comparable ==============")
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)

	fmt.Println("find result", v.Get(v.Find(5)))

	// un-comparable type
	fmt.Println("\n=========== un-comparable ==============")
	var vv = New([]UnComparable{{Key: 1, Value: []int{11}}, {Key: 3, Value: []int{33}}, {Key: 2, Value: []int{22}}, {Key: 4, Value: []int{44}}})
	fmt.Println("find result", vv.Find(UnComparable{Key: 2}))
}

func TestVector_FindBy(t *testing.T) {
	// comparable type
	fmt.Println("=========== comparable ==============")
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)

	fmt.Println("find result", v.Get(v.FindBy(func(x int) bool {
		return x == 5
	})))

	// un-comparable type
	fmt.Println("\n=========== un-comparable ==============")
	var vv = New([]UnComparable{{Key: 1, Value: []int{11}}, {Key: 3, Value: []int{33}}, {Key: 2, Value: []int{22}}, {Key: 4, Value: []int{44}}})
	fmt.Println("find result", vv.Get(vv.FindBy(func(x UnComparable) bool {
		return x.Key == 2
	})))
}

func TestVector_Erase(t *testing.T) {
	// comparable type
	fmt.Println("=========== comparable ==============")
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)
	v.Sort(func(a, b int) bool {
		return a < b
	})
	fmt.Println("before", v)
	v.EraseAt(v.Find(5))
	fmt.Println("after", v)

	// un-comparable type
	fmt.Println("\n=========== un-comparable ==============")
	var vv = New([]UnComparable{{Key: 1, Value: []int{11}}, {Key: 3, Value: []int{33}}, {Key: 2, Value: []int{22}}, {Key: 4, Value: []int{44}}})
	fmt.Println("before", vv)
	vv.EraseAt(vv.Find(UnComparable{Key: 2}))
	fmt.Println("after", vv)
}

func TestVector_Begin(t *testing.T) {
	// comparable type
	fmt.Println("=========== comparable ==============")
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)
	it := v.Begin()
	for it.Next() {
		fmt.Printf("%v ", it.Value())
	}

	fmt.Println()

	it = v.RBegin()
	for it.Next() {
		fmt.Printf("%v ", it.Value())
	}

	fmt.Println()

	// un-comparable type
	fmt.Println("\n=========== un-comparable ==============")
	var vv = New([]UnComparable{{Key: 1, Value: []int{11}}, {Key: 3, Value: []int{33}}, {Key: 2, Value: []int{22}}, {Key: 4, Value: []int{44}}})
	it2 := vv.Begin()
	for it2.Next() {
		fmt.Println(it2.Value())
	}
}

func TestVector_Insert(t *testing.T) {
	// comparable type
	fmt.Println("=========== comparable ==============")
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)
	v.Sort(func(a, b int) bool {
		return a < b
	})
	fmt.Println("before", v.Slices())
	v.Insert(7, 10, 11, 12, 13, 14, 15, 16, 17)
	fmt.Println("after", v.Slices())
}

func TestVector_Join(t *testing.T) {
	// comparable type
	fmt.Println("=========== comparable ==============")
	var intArray = []int{1, 3, 2, 5, 4, 7, 6, 9, 8}
	var v = New(intArray)
	fmt.Println("join", v.Join("|"))
}

func TestVector_MarshalJSON(t *testing.T) {
	var vec = New([]int{})
	vec.Append(1, 2, 3, 4, 5)
	jsonStr, err := json.Marshal(vec)
	fmt.Println(string(jsonStr), err)

	// unmarshal
	var vec2 = New[int](nil)
	err = json.Unmarshal(jsonStr, &vec2)
	fmt.Println(vec2.Slices(), err)

	jsonStr, err = json.Marshal(3)
	fmt.Println(string(jsonStr), err)
}
