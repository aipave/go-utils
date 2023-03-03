package heap

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	var hp = New[int]()
	var arr = []int{3, 1, 5, 8, 23, 22, 12, 89, 100, 44, 2, 7, 9, 13}
	hp.Push(arr...)

	for hp.Size() > 0 {
		fmt.Println(hp.Pop())
	}
}
