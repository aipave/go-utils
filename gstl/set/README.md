## Set and MultiSet

### Requirements
Go version >= 1.18.1

### Usage

```go
go get -u gitlab-vywrajy.micoworld.net/yoho-go/container
```

import package
```go
import "gitlab-vywrajy.micoworld.net/yoho-go/container/set"
```

1. Create
> NewSet: only support int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64,float32,float64,string类型
> NewExtSet: need Less func
> NewMultiSet: same to NewSet
> NewMultiExtSet: same to NewMultiSet
```go
type User struct {
    UID int64
    Age int64
    Avatars []string
}

func less(a, b User) int {
    return a.uid - b.uid
}

func handler() {
    // 1. Create
    var s = set.NewSet[int]()

    // 2. Create with lock
    var s2 = set.NewSet[int](options.WithLock())

    // 3. Create Ext Set
    var s3 = set.NewExtSet[User](less)

    // 4. Create Multi Set
    var s4 = set.NewMultiSet[int]()

    // 5. Create Multi Set With Lock
    var s5 = set.NewMultiSet[int](options.WithLock())
}
```

2. Insert and Erase
```go
func handler() {
    var list = []int{1, 3, 2, 4, 5, 7, 8, 9}
    var s = set.NewSet[int]()

    // Insert
    s.Insert(list...) // 1, 2, 3, 4, 5, 6, 7, 8, 9

    // Erase
    s.Erase(5) // 1, 2, 3, 4, 6, 7, 8, 9
}
```

3. Slices, Size and Clear
```go
func handler() {
    var list = []int{1, 3, 2, 4, 5, 7, 8, 9}
    var s = set.NewSet[int]()

    // Insert
    s.Insert(list...) // 1, 2, 3, 4, 5, 6, 7, 8, 9

    fmt.Println("size", s.Size())
    fmt.Println("slices", s.Slices()) // [1,2,3,4,5,6,7,8,9]
    s.Clear()
    fmt.Println("size", s.Size())
}
```

4. Iterator
```go
func handler() {
   var list = []int{1, 3, 2, 4, 5, 7, 8, 9}
   var s = set.NewSet[int]()

   // Insert
   s.Insert(list...) // 1, 2, 3, 4, 5, 6, 7, 8, 9

   // iterator from left to right
   it := s.Begin()
   for it.Next() {
       fmt.Println("%v ", it.Value())
   }
   fmt.Println()

   // iterator from right to left
   it := s.RBegin()
   for it.Next() {
       fmt.Println("%v ", it.Value())
   }
   fmt.Println()

   // iterator from find position to right
   it := s.BeginFrom(5)
   for it.Next() {
       fmt.Println("%v ", it.Value())
   }
}
```

5. Intersect, Union and Diff
```go
func handler() {
	var a, b = NewSet[int](), NewSet[int]()

	for i := 1; i < 20; i++ {
		a.Insert(i)
	}

	for i := 16; i < 36; i++ {
		b.Insert(i)
	}

	fmt.Println("intersect", a.Intersect(b).Slices())
	fmt.Println("union", a.Union(b).Slices())
	fmt.Println("diff", a.Diff(b).Slices())
}
```
