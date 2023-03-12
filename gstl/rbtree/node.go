package rbtree

type Color bool

const (
	Red   = false
	Black = true
)

type Node[K any, V any] struct {
	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
	color  Color
	key    K
	value  V
}

func (n *Node[K, V]) Key() K {
	return n.key
}

func (n *Node[K, V]) Value() V {
	return n.value
}

func (n *Node[K, V]) SetValue(val V) {
	n.value = val
}

func (n *Node[K, V]) Next() *Node[K, V] {
	return next(n)
}

func next[K any, V any](ptr *Node[K, V]) *Node[K, V] {
	if ptr.right != nil {
		return leftNode(ptr.right)
	}

	parent := ptr.parent
	for parent != nil && ptr == parent.right {
		ptr, parent = parent, parent.parent
	}
	return parent
}

func (n *Node[K, V]) Prev() *Node[K, V] {
	return prev(n)
}

func prev[K any, V any](ptr *Node[K, V]) *Node[K, V] {
	if ptr.left != nil {
		return rightNode(ptr.left)
	}

	parent := ptr.parent
	for parent != nil && parent.left == ptr {
		ptr, parent = parent, parent.parent
	}
	return parent
}

// 获取一个节点的最左节点
func leftNode[K any, V any](x *Node[K, V]) *Node[K, V] {
	for x.left != nil {
		x = x.left
	}
	return x
}

// 获取一个节点的最右节点
func rightNode[K any, V any](x *Node[K, V]) *Node[K, V] {
	for x.right != nil {
		x = x.right
	}
	return x
}
