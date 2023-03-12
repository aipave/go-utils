package fns

import (
	"fmt"
	"testing"
)

func TestIf(t *testing.T) {
	fmt.Println(If(true, 1, 2))
	fmt.Println(If(false, 1, 2))
}
