package main

import (
	"cmp"
	"fmt"
)

// Properties:
// - Root is black
// - Red nodes have only black children
//
// Insertion:
// - Regular insert
// - New node is red
// - Fix properties
//
// Fix properties:
//   Scenario 1: Inserted node is root --> color black
//   Scenario 2: Uncle is red --> recolor parent, grandparent and uncle
//   Scenario 3: Uncle is black (trinagle) --> rotate parent in opposite direction
//   Scenario 4: Uncle is black (line) --> rotate grandparent in opposite direction and recolor original parent and grandparent
//
// Json Tree Visualizer: https://vanya.jp.net/vtree/
//

type Color int

const (
	Black Color = 0
	Red   Color = 1
)

type Relationship int

const (
	Triangle Relationship = 0
	Line     Relationship = 1
)

func (c Color) String() string {
	if c == Black {
		return "black"
	} else {
		return "red"
	}
}

type Value interface {
	cmp.Ordered
}

func String[V Value](v V) string {
	return fmt.Sprintf("%v", v)
}

type Tree[V Value] struct {
	node *Node[V]
}

func MakeTree[V Value]() Tree[V] {
	return Tree[V]{}
}

type Node[V Value] struct {
	val   V
	color Color
	p     *Node[V]
	left  *Node[V]
	right *Node[V]
}

func (t *Tree[V]) Insert(v V) {
	if t.node == nil {
		t.node = &Node[V]{val: v, color: Black}
	} else {
		t.node.Insert(v)
	}
	t.node.color = Black
}

func (t Tree[V]) String() string {
	if t.node == nil {
		return "EmptyTree"
	} else {
		return t.node.String()
	}
}

func (t Tree[V]) Json() string {
	if t.node == nil {
		return "{}"
	} else {
		return t.node.Json()
	}
}

func (t Tree[V]) Contains(v V) bool {
	return t.node.Contains(v)
}

func (t Tree[V]) Height() int {
	return t.node.Height(0)
}

func (t Tree[V]) Size() int {
	return t.node.Size()
}

func (n *Node[V]) Insert(v V) {
	if v == n.val {
		return
	} else if v > n.val {
		if n.right == nil {
			n.right = &Node[V]{val: v, color: Red, p: n}
			n.right.FixViolations()
		} else {
			n.right.Insert(v)
		}
	} else {
		if n.left == nil {
			n.left = &Node[V]{val: v, color: Red, p: n}
			n.left.FixViolations()
		} else {
			n.left.Insert(v)
		}
	}
}

func (n *Node[V]) FixViolations() {
	if n.p == nil || n.p.p == nil || n.p.color != Red {
		return
	}

	uncle, rel, ok := n.Uncle()

	// leaves are black, so no uncle means black
	redUncle := ok && uncle.color == Red

	if redUncle {
		n.p.Recolor()
		n.p.p.Recolor()
		uncle.color = Black
		n.p.p.FixViolations()
	} else {
		var rotateNode *Node[V]
		if rel == Triangle {
			rotateNode = n.p
		} else {
			rotateNode = n.p.p
			n.p.Recolor()
			n.p.p.Recolor()
		}

		if n.p.left == n {
			rotateNode.RightRotate()
			rotateNode.right.FixViolations()
		} else {
			rotateNode.LeftRotate()
			rotateNode.left.FixViolations()
		}
	}
}

func (n *Node[V]) Recolor() {
	if n.color == Black {
		n.color = Red
	} else {
		n.color = Black
	}
}

func (n *Node[V]) LeftRotate() {
	if n.right == nil {
		panic("Can't left-rotate if I don't have a right child")
	}

	var newLeftChild Node[V] = *n // store old n, becomes left child
	*n = *n.right                 // right child becomes the parent

	n.p = newLeftChild.p
	newLeftChild.p = n

	if n.left == nil {
		newLeftChild.right = nil
	} else {
		*newLeftChild.right = *n.left
		if newLeftChild.right.left != nil {
			newLeftChild.right.left.p = newLeftChild.right
		}
		if newLeftChild.right.right != nil {
			newLeftChild.right.right.p = newLeftChild.right
		}
		newLeftChild.right.p = &newLeftChild
	}
	n.left = &newLeftChild
	if n.right != nil {
		n.right.p = n
	}

	if newLeftChild.left != nil {
		newLeftChild.left.p = &newLeftChild
	}
}

func (n *Node[V]) RightRotate() {
	if n.left == nil {
		panic("Can't right-rotate if I don't have a left child")
	}

	var newRightChild Node[V] = *n // store old n, becomes right child
	*n = *n.left                   // left child becomes the parent

	n.p = newRightChild.p
	newRightChild.p = n

	if n.right == nil {
		newRightChild.left = nil
	} else {
		*newRightChild.left = *n.right
		if newRightChild.left.left != nil {
			newRightChild.left.left.p = newRightChild.left
		}
		if newRightChild.left.right != nil {
			newRightChild.left.right.p = newRightChild.left
		}
		newRightChild.left.p = &newRightChild
	}
	n.right = &newRightChild
	if n.left != nil {
		n.left.p = n
	}

	if newRightChild.right != nil {
		newRightChild.right.p = &newRightChild
	}
}

func (n *Node[V]) Contains(v V) bool {
	if n == nil {
		return false
	}

	if n.val == v {
		return true
	}

	if v > n.val {
		return n.right.Contains(v)
	} else {
		return n.left.Contains(v)
	}

}

func (n *Node[V]) String() string {
	if n == nil {
		return ""
	}

	acc := "L = {"
	acc += n.left.String()
	acc += "} "
	acc += String(n.val)
	acc += " R = {"
	acc += n.right.String()
	acc += "}"

	return acc
}

func (n *Node[V]) Json() string {
	if n == nil {
		return "{}"
	}

	var parent string
	if n.p == nil {
		parent = "\"nil\""
	} else {
		parent = String(n.p.val)
		// parent = fmt.Sprintf("\"%p\"", n.p)
	}

	acc := "{"
	acc += "\"value\": " + String(n.val) + ","
	// acc += "\"ref\": \"" + fmt.Sprintf("%p", n) + "\","
	acc += "\"color\": \"" + n.color.String() + "\","
	acc += "\"parent\": " + parent + ","
	acc += "\"left\": " + n.left.Json() + ","
	acc += "\"right\": " + n.right.Json()
	acc += "}"

	return acc
}

func (n *Node[V]) Height(depth int) int {
	if n == nil {
		return depth
	}
	l := n.left.Height(depth + 1)
	r := n.right.Height(depth + 1)

	if l > r {
		return l
	} else {
		return r
	}
}

func (n *Node[V]) Size() int {
	if n == nil {
		return 0
	}

	l := n.left.Size()
	r := n.right.Size()

	return l + r + 1
}

func (n *Node[V]) Uncle() (*Node[V], Relationship, bool) {
	parent := n.p
	if parent == nil {
		return nil, 0, false
	}

	grandparent := parent.p
	if grandparent == nil {
		return nil, 0, false
	}

	leftChild := parent.left == n
	parentIsLeftChild := grandparent.left == parent
	var rel Relationship
	if leftChild == parentIsLeftChild {
		rel = Line
	} else {
		rel = Triangle
	}

	var uncle *Node[V]
	if grandparent.left == parent {
		uncle = grandparent.right
	} else {
		uncle = grandparent.left
	}
	return uncle, rel, uncle != nil
}
