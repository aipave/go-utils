package rbtree

import "github.com/alyu01/go-utils/gstl/compare"

type RBTree[K any, V any] struct {
    root   *Node[K, V]
    size   int
    keyCmp compare.Less[K]
}

func New[K any, V any](keyCmp compare.Less[K]) *RBTree[K, V] {
    return &RBTree[K, V]{keyCmp: keyCmp}
}

func (rbt *RBTree[K, V]) Clear() {
    rbt.root = nil
    rbt.size = 0
}

func (rbt *RBTree[K, V]) leftRotate(x *Node[K, V]) {
    y := x.right
    x.right = y.left
    if y.left != nil {
        y.left.parent = x
    }
    y.parent = x.parent
    if x.parent == nil {
        rbt.root = y
    } else if x == x.parent.left {
        x.parent.left = y
    } else {
        x.parent.right = y
    }
    y.left = x
    x.parent = y
}

func (rbt *RBTree[K, V]) rightRotate(x *Node[K, V]) {
    y := x.left
    x.left = y.right
    if y.right != nil {
        y.right.parent = x
    }
    y.parent = x.parent
    if x.parent == nil {
        rbt.root = y
    } else if x == x.parent.right {
        x.parent.right = y
    } else {
        x.parent.left = y
    }
    y.right = x
    x.parent = y
}

func (rbt *RBTree[K, V]) FindNode(key K) *Node[K, V] {
    ptr := rbt.root
    for ptr != nil {
        if rbt.keyCmp(key, ptr.key) < 0 {
            ptr = ptr.left
        } else if rbt.keyCmp(key, ptr.key) > 0 {
            ptr = ptr.right
        } else {
            return ptr
        }
    }

    return nil
}

func (rbt *RBTree[K, V]) FindFirstNode(key K) (first *Node[K, V]) {
    first = rbt.FindNode(key)
    if first == nil {
        return
    }

    for prev := first.Prev(); prev != nil && rbt.keyCmp(key, prev.key) == 0; prev = prev.Prev() {
        first = prev
    }
    return
}

func (rbt *RBTree[K, V]) FindLastNode(key K) (last *Node[K, V]) {
    last = rbt.FindNode(key)
    if last == nil {
        return
    }

    for next := last.Next(); next != nil && rbt.keyCmp(key, next.key) == 0; next = next.Next() {
        last = next
    }
    return
}

func (rbt *RBTree[K, V]) Size() int {
    return rbt.size
}

func (rbt *RBTree[K, V]) Insert(key K, value V) {
    ptr := rbt.root

    var parent *Node[K, V]
    for ptr != nil {
        parent = ptr
        if rbt.keyCmp(key, ptr.key) < 0 {
            ptr = ptr.left
        } else {
            ptr = ptr.right
        }
    }

    z := &Node[K, V]{parent: parent, color: Red, key: key, value: value}
    rbt.size++

    if parent == nil { // 空树
        z.color = Black
        rbt.root = z
        return

    } else if rbt.keyCmp(z.key, parent.Key()) < 0 { // 插入到目标节点左边
        parent.left = z

    } else {
        parent.right = z // 插入到目标节点右边
    }

    rbt.insertBalance(z)
}

func (rbt *RBTree[K, V]) insertBalance(ptr *Node[K, V]) {
    // 父节点不为空 && 父节点为红色 递归调整
    for ptr.parent != nil && ptr.parent.color == Red {
        var parent = ptr.parent
        var grandParent = parent.parent // 由于父节点为红色，所以父节点的父节点一定不为空 （红黑树特效保证根节点一定为黑色）

        if parent == grandParent.left {
            uncle := grandParent.right
            if uncle != nil && uncle.color == Red {
                parent.color = Black
                uncle.color = Black
                grandParent.color = Red
                ptr = grandParent

            } else {
                if ptr == parent.right {
                    ptr = ptr.parent
                    rbt.leftRotate(ptr)
                }
                ptr.parent.color = Black
                ptr.parent.parent.color = Red
                rbt.rightRotate(ptr.parent.parent)
            }

        } else {
            uncle := grandParent.left
            if uncle != nil && uncle.color == Red {
                parent.color = Black
                uncle.color = Black
                grandParent.color = Red
                ptr = grandParent

            } else {
                if ptr == parent.left {
                    ptr = ptr.parent
                    rbt.rightRotate(ptr)
                }
                ptr.parent.color = Black
                ptr.parent.parent.color = Red
                rbt.leftRotate(ptr.parent.parent)
            }
        }
    }

    rbt.root.color = Black
}

func (rbt *RBTree[K, V]) Delete(node *Node[K, V]) {
    ptr := node
    if ptr == nil {
        return
    }

    var x, y *Node[K, V]
    if ptr.left != nil && ptr.right != nil {
        y = next(ptr)
    } else {
        y = ptr
    }

    if y.left != nil {
        x = y.left
    } else {
        x = y.right
    }

    xParent := y.parent
    if x != nil {
        x.parent = xParent
    }

    if y.parent == nil {
        rbt.root = x
    } else if y == y.parent.left {
        y.parent.left = x
    } else {
        y.parent.right = x
    }

    if y != ptr {
        ptr.key = y.key
        ptr.value = y.value
    }

    if y.color == Black {
        rbt.deleteBalance(x, xParent)
    }
    rbt.size--
}

func (rbt *RBTree[K, V]) deleteBalance(x, parent *Node[K, V]) {
    var w *Node[K, V]
    for x != rbt.root && getColor(x) == Black {
        if x != nil {
            parent = x.parent
        }
        if x == parent.left {
            x, w = rbt.leftBalance(x, parent, w)
        } else {
            x, w = rbt.rightBalance(x, parent, w)
        }
    }

    if x != nil {
        x.color = Black
    }
}

func (rbt *RBTree[K, V]) leftBalance(x, parent, w *Node[K, V]) (*Node[K, V], *Node[K, V]) {
    w = parent.right
    if w.color == Red {
        w.color = Black
        parent.color = Red
        rbt.leftRotate(parent)
        w = parent.right
    }

    if getColor(w.left) == Black && getColor(w.right) == Black {
        w.color = Red
        x = parent
    } else {
        if getColor(w.right) == Black {
            if w.left != nil {
                w.left.color = Black
            }
            w.color = Red
            rbt.rightRotate(w)
            w = parent.right
        }
        w.color = parent.color
        parent.color = Black
        if w.right != nil {
            w.right.color = Black
        }
        rbt.leftRotate(parent)
        x = rbt.root
    }
    return x, w
}

func (rbt *RBTree[K, V]) rightBalance(x, parent, w *Node[K, V]) (*Node[K, V], *Node[K, V]) {
    w = parent.left
    if w.color == Red {
        w.color = Black
        parent.color = Red
        rbt.rightRotate(parent)
        w = parent.left
    }

    if getColor(w.left) == Black && getColor(w.right) == Black {
        w.color = Red
        x = parent
    } else {
        if getColor(w.left) == Black {
            if w.right != nil {
                w.right.color = Black
            }
            w.color = Red
            rbt.leftRotate(w)
            w = parent.left
        }
        w.color = parent.color
        parent.color = Black
        if w.left != nil {
            w.left.color = Black
        }
        rbt.rightRotate(parent)
        x = rbt.root
    }
    return x, w
}

func getColor[K any, V any](n *Node[K, V]) Color {
    if n == nil {
        return Black
    }

    return n.color
}

func (rbt *RBTree[K, V]) Begin() *Node[K, V] {
    if rbt.root == nil {
        return nil
    }
    return leftNode(rbt.root)
}

func (rbt *RBTree[K, V]) RBegin() *Node[K, V] {
    if rbt.root == nil {
        return nil
    }
    return rightNode(rbt.root)
}
