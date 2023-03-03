package test_example

import (
    "testing"

    "github.com/alyu01/go-utils/gstl/vector"
)

func TestGSTL(t *testing.T) {
    s := []int{1, 2, 3}
    idx := vector.New(s).Find(1)
    if idx < 0 {
        t.Fatal("not exits")
    }
    t.Logf("exits")

}
