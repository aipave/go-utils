package rbtree

import (
	"fmt"
	"testing"

	"github.com/aipave/go-utils/gstl/compare"
)

type n struct {
	K int
	V int
}

func TestNew(t *testing.T) {
	var tree = New[int, n](compare.OrderedLess[int])
	for i := 1; i <= 100; i++ {
		tree.Insert(i%10, n{K: i % 10, V: i})
	}

	for node := tree.Begin(); node != nil; node = node.Next() {
		fmt.Printf("%v ", node.Value())
	}

	fmt.Println()

	for node := tree.RBegin(); node != nil; node = node.Prev() {
		fmt.Printf("%v ", node.value)
	}

	fmt.Println()

	fmt.Println(tree.FindNode(1).Value(), tree.FindFirstNode(1).Value(), tree.FindLastNode(1).Value())
	fmt.Println(tree.FindNode(2).Value(), tree.FindFirstNode(2).Value(), tree.FindLastNode(2).Value())
}
