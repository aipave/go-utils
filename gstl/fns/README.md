## fns

generic template function

### Requirements
Go version >= 1.18.1

### Usage

1. If Ternary operator
```
func main() {
    var a = fns.If(true, 1, 2)
    var b = fns.If(false, 1, 2)

    fmt.Println("a = ", a, ", b = ", b) // a = 1 , b = 2
}
```
