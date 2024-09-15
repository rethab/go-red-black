// A sample package that implements a [RedBlack] tree in Go.
//
// This is a learning module and should not be used in real programs.
//
// (Some of the) Properties:
//   - Root is black
//   - Red nodes have only black children
//
// Insertion:
//   - Regular insert
//   - New node is red
//   - Fix properties
//
// Fix properties:
//   - Scenario 1: Inserted node is root --> color black
//   - Scenario 2: Uncle is red --> recolor parent, grandparent and uncle
//   - Scenario 3: Uncle is black (trinagle) --> rotate parent in opposite direction
//   - Scenario 4: Uncle is black (line) --> rotate grandparent in opposite direction and recolor original parent and grandparent
//
// [RedBlack]: https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
package redblack

import (
	"cmp"
	"fmt"
)

type color int

const (
	black color = 0
	red   color = 1
)

type relationship int

const (
	triangle relationship = 0
	line     relationship = 1
)

func (c color) String() string {
	if c == black {
		return "black"
	} else {
		return "red"
	}
}

// The type constraint that values in a tree need to satisfy
type Value interface {
	cmp.Ordered
}

func show[V Value](v V) string {
	return fmt.Sprintf("%v", v)
}

type Tree[V Value] struct {
	node *node[V]
}

// MakeTree creates a new Red-Black Tree
func MakeTree[V Value]() Tree[V] {
	return Tree[V]{}
}

type node[V Value] struct {
	val   V
	color color
	p     *node[V]
	left  *node[V]
	right *node[V]
}

// Insert a value into a the tree.
//
// If the value already exists, nothing happens
func (t *Tree[V]) Insert(v V) {
	if t.node == nil {
		t.node = &node[V]{val: v, color: black}
	} else {
		t.node.insert(v)
	}
	t.node.color = black
}

// Formats the string in a human readable format
func (t Tree[V]) String() string {
	if t.node == nil {
		return "EmptyTree"
	} else {
		return t.node.String()
	}
}

// Format the tree with JSON
//
// This can be put into a [JSONVisualizer] for debug purposes
//
// [JSONVisualizer]: https://vanya.jp.net/vtree/
func (t Tree[V]) Json() string {
	if t.node == nil {
		return "{}"
	} else {
		return t.node.json()
	}
}

// Checks whether the specified value is in the tree
func (t Tree[V]) Contains(v V) bool {
	return t.node.contains(v)
}

// Returns the height of the tree
//
// The height is the number of nodes from the root to the
// furthest-away leaf.
//
// Since a red-black tree is self-balancing, this is an important
// property. A naive binary search tree could potentially have
// height O(n), whereas as self-balancing binary search tree should
// approximate O(log n)
func (t Tree[V]) Height() int {
	return t.node.height(0)
}

// Returns the total number of nodes in the tree
func (t Tree[V]) Size() int {
	return t.node.size()
}

func (n *node[V]) insert(v V) {
	if v == n.val {
		return
	} else if v > n.val {
		if n.right == nil {
			n.right = &node[V]{val: v, color: red, p: n}
			n.right.fixViolations()
		} else {
			n.right.insert(v)
		}
	} else {
		if n.left == nil {
			n.left = &node[V]{val: v, color: red, p: n}
			n.left.fixViolations()
		} else {
			n.left.insert(v)
		}
	}
}

func (n *node[V]) fixViolations() {
	if n.p == nil || n.p.p == nil || n.p.color != red {
		return
	}

	uncle, rel, ok := n.uncle()

	// leaves are black, so no uncle means black
	redUncle := ok && uncle.color == red

	if redUncle {
		n.p.recolor()
		n.p.p.recolor()
		uncle.color = black
		n.p.p.fixViolations()
	} else {
		var rotateNode *node[V]
		if rel == triangle {
			rotateNode = n.p
		} else {
			rotateNode = n.p.p
			n.p.recolor()
			n.p.p.recolor()
		}

		if n.p.left == n {
			rotateNode.rightRotate()
			rotateNode.right.fixViolations()
		} else {
			rotateNode.leftRotate()
			rotateNode.left.fixViolations()
		}
	}
}

func (n *node[V]) recolor() {
	if n.color == black {
		n.color = red
	} else {
		n.color = black
	}
}

func (n *node[V]) leftRotate() {
	if n.right == nil {
		panic("Can't left-rotate if I don't have a right child")
	}

	var newLeftChild node[V] = *n // store old n, becomes left child
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

func (n *node[V]) rightRotate() {
	if n.left == nil {
		panic("Can't right-rotate if I don't have a left child")
	}

	var newRightChild node[V] = *n // store old n, becomes right child
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

func (n *node[V]) contains(v V) bool {
	if n == nil {
		return false
	}

	if n.val == v {
		return true
	}

	if v > n.val {
		return n.right.contains(v)
	} else {
		return n.left.contains(v)
	}

}

func (n *node[V]) String() string {
	if n == nil {
		return ""
	}

	acc := "L = {"
	acc += n.left.String()
	acc += "} "
	acc += show(n.val)
	acc += " R = {"
	acc += n.right.String()
	acc += "}"

	return acc
}

func (n *node[V]) json() string {
	if n == nil {
		return "{}"
	}

	var parent string
	if n.p == nil {
		parent = "\"nil\""
	} else {
		parent = show(n.p.val)
		// parent = fmt.Sprintf("\"%p\"", n.p)
	}

	acc := "{"
	acc += "\"value\": " + show(n.val) + ","
	// acc += "\"ref\": \"" + fmt.Sprintf("%p", n) + "\","
	acc += "\"color\": \"" + n.color.String() + "\","
	acc += "\"parent\": " + parent + ","
	acc += "\"left\": " + n.left.json() + ","
	acc += "\"right\": " + n.right.json()
	acc += "}"

	return acc
}

func (n *node[V]) height(depth int) int {
	if n == nil {
		return depth
	}
	l := n.left.height(depth + 1)
	r := n.right.height(depth + 1)

	if l > r {
		return l
	} else {
		return r
	}
}

func (n *node[V]) size() int {
	if n == nil {
		return 0
	}

	l := n.left.size()
	r := n.right.size()

	return l + r + 1
}

func (n *node[V]) uncle() (*node[V], relationship, bool) {
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
	var rel relationship
	if leftChild == parentIsLeftChild {
		rel = line
	} else {
		rel = triangle
	}

	var uncle *node[V]
	if grandparent.left == parent {
		uncle = grandparent.right
	} else {
		uncle = grandparent.left
	}
	return uncle, rel, uncle != nil
}
