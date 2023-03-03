## Vector

### Requirements
Go version >= 1.18.1

### Usage

```go
go get -u gitlab-vywrajy.micoworld.net/yoho-go/container
```

import package
```go
import "gitlab-vywrajy.micoworld.net/yoho-go/container/vector"
```

1. Create a Vector
```go
type User struct {
    UID int64
    Age int64
    Avatars []string
}

func handler() {
    // 1. Create an int Vector without lock
    var vec = vector.New[int](nil)

    var list []int
    var vec2 = vector.New[int](list) // ok

    // 2. Create an Struct Vector
    var structVec = vector.New[User](nil)

    var ary []User
    var structVec3 = vector.New(ary) // ok
}
```

2. Sort and Reverse
```go
func sort() {
    var vec = vector.New[int]([1, 3, 2, 4, 6, 5, 8, 7, 10, 9]]) // create an unsafeVec

    // sort
    vec.Sort(func a, b int) bool {
        return a < b // asc
    })
    fmt.Println("sorted", vec.Slices()) // result: 1 2 3 4 5 6 7 8 9 10

    // reverse
    vec.Reverse()
    fmt.Println("reverse", vec.Slices()) // result: 10 9 8 7 6 5 4 3 2 1
}
```

3. Find and FindBy
```go
var User struct {
    Name string
    Age  int
}

func find() {
    // find
    var vec = vector.New[int]([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]) // type int

    pos := vec.Find(4) // -1 means not found
    fmt.Println("found pos: ", pos, "pos value: ", vec.Get(pos))

    var uVec = vector.New[User](nil)
    var u = User{Name: "yangmingzi", Age: 28}
    uVec.Append(u)

    filter := User{Name: "yangmingzi"} // filter by name is not allowed
    pos = uVec.Find(filter) // not found
    fmt.Println("found pos: ", pos)
    fmt.Println("found pos: ", uVec.Find(u)) // ok

    // findBy
    pos = uVec.FindBy(func (x User) bool {
        return x.Name == "yangmingzi" // filter by name
    })
    fmt.Println("found pos: ", pos, "found value: ", uVec.Get(pos)) // ok
}
```

4. Erase
```go
func erase() {
    // 1. comparable type
    var vec = vector.New([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]) // type int

    vec.Erase(4) // ok, remove value 4
    vec.EraseAt(4) // ok, remove index 4
}
```

5. Join
```go
func join() {
    var vec = vector.New([1, 2, 3, 4, 5, 6])

    fmt.Println(vec.Join("|")) // result: 1|2|3|4|5|6
    fmt.Println(vec.Join(",")) // result: 1,2,3,4,5,6
}
```

6. Iterator and ReverseIterator
```go
func iterator() {
    var vec = vector.New([1, 2, 3, 4, 5, 6])

    iter := vec.Begin()
    for iter.Next() {
        fmt.Println("value: ", iter.Value())
    }

    iter = vec.RBegin()
    for iter.Next() {
        fmt.Println("value: ", iter.Value())
    }
}
```
