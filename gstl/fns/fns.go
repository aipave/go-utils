package fns

// If Ternary operator
func If[T any](isA bool, a T, b T) T {
    if isA {
        return a
    }

    return b
}

type Pair[K any, V any] struct {
    Key   K `json:"key"`
    Value V `json:"value"`
}
