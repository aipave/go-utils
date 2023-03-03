## Map and MultiMap

### Requirements
Go version >= 1.18.1

### Usage

```go
```

import package
```go
```

1. Create a Map or MultiMap
```go
type User struct {
    UID int64
    Age int64
    Avatars []string
}

func less(a, b User) int {
    return a.UID - b.UID
}

func handler() {
    // 1. Create an int Vector without lock
    var Map = smap.New[int, int]()

    var Map2 = smap.New[int, User]()

    // 2. Create an Struct Vector
    var structMap = smap.New[int, User](options.WithLocker())

    var structMap3 = smap.New[int, int](options.WithLocker()) // ok

    // 3. Create Struct key Map
    var defMap = smap.NewExt[User, bool](less)

    var defMap2 = smap.NewExt[User, bool](less, options.WithLocker())

    // 3. Create multi-map
    var multiMap = smap.NewMultiMap[int, int]()

    var multiMap2 = smap.NewMultiExtMap[User, bool]()
}
```

2. Get and Set
```
func main() {
    var m = smap.New[int, int]()

    fmt.Println("set", m.Set(1, 1)) // return ok if key does't exist
    fmt.Println("set", m.Set(1, 2)) // return false means key is exist, then update value to 2

    if v, ok := m.GetOk(1); ok {
        fmt.Println("get", v)
    }

    if v, ok := m.GetOk(2); ok {
        fmt.Println("get", v)
    }

    fmt.Println("get", m.Get(2)) // return default value(e.g. 0, "", ...) if key does't exist

	fmt.Println("get set", m.GetSet(1, 3)) // get preview value and set new value
	fmt.Println("get set", m.GetSet(1, 4))

    var mm = smap.NewMultiMap[int, int]()

    mm.Set(1, 2)
    mm.Set(1, 1)
    mm.Set(1, 3)

    fmt.Println("get", mm.Get(1)) // it will return then first element visited in the tree

    it := mm.GetFirst(1)
    for it.Next() {
        fmt.Println("get", it.Value()) // from left to right, 2 1 3
    }

    it = mm.GetLast(1)
    for it.Next() {
        fmt.Println("get", it.Value()) // from right to left, 3 1 2
    }
}
```

3. Keys and Values
```
func keysAndValues() {
	var m = NewMap[int, int]()

	for i := 1; i <= 10; i++ {
		m.Set(i, i*2)
	}

	fmt.Println("keys", m.Keys()) // 1 2 3 4 5 6 7 8 9 10
	fmt.Println("values", m.Values()) // 2 4 6 8 10 12 14 16 18 20
}

```

4. Contains
```

func contains() {
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
```

5. Iterator
```
func iterator() {
	var m = NewMap[int, int]()

	for i := 1; i <= 10; i++ {
		m.Set(i, i)
	}

	it := m.Begin()
	for it.Next() {
		fmt.Println("iter", it.Value())
	}

	it = m.RBegin()
	for it.Next() {
		fmt.Println("iter", it.Value())
	}

	it = m.BeginFrom(3)
	for it.Next() {
		fmt.Println("iter", it.Value())
	}

	it = m.RBeginFrom(3)
	for it.Next() {
		fmt.Println("iter", it.Value())
	}
}
```