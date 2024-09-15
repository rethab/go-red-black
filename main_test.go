package main

import (
	"cmp"
	"fmt"
	"testing"
)

func TestEmptyNotContains(t *testing.T) {
	tree := MakeTree[int]()
	if tree.Contains(5) {
		t.Fatalf("Empty tree contains 5")
	}
	if tree.Size() != 0 {
		t.Fatalf("Empty tree has not size 0")
	}
	if tree.Height() != 0 {
		t.Fatalf("Empty tree has not height 0")
	}
}

func TestString(t *testing.T) {
	tree := MakeTree[string]()
	tree.Insert("hello")
	if tree.Size() != 1 {
		t.Fatalf("tree has not size 1")
	}
	tree.Insert("world")
	if tree.Size() != 2 {
		t.Fatalf("tree has not size 2 after inserting second word")
	}
	if !tree.Contains("world") {
		t.Fatalf("tree does not contain 'world' anymore")
	}
}

func TestFloat(t *testing.T) {
	tree := MakeTree[float64]()
	tree.Insert(1.2)
	if tree.Size() != 1 {
		t.Fatalf("tree has not size 1")
	}
	tree.Insert(4.3)
	if tree.Size() != 2 {
		t.Fatalf("tree has not size 2 after inserting second value")
	}
	if !tree.Contains(1.2) {
		t.Fatalf("tree does not contain 1.2 anymore")
	}
}

func TestInsertDuplicate(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(5)
	if tree.Size() != 1 {
		t.Fatalf("tree has not size 1")
	}
	tree.Insert(5)
	if tree.Size() != 1 {
		t.Fatalf("tree has not size 1 after inserting 5 twice")
	}
	if !tree.Contains(5) {
		t.Fatalf("tree does not contain 5 anymore")
	}
}

func TestInsertContainsOne(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(5)
	if !tree.Contains(5) {
		t.Fatalf("Tree does not contain 5")
	}
	if tree.Size() != 1 {
		t.Fatalf("Tree has not size 1")
	}
	if tree.Height() != 1 {
		t.Fatalf("Tree has not height 1")
	}
}

func TestInsertContainsLeft(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(5)
	tree.Insert(1)
	if !tree.Contains(1) {
		t.Fatalf("Tree does not contain 1")
	}
	if tree.Size() != 2 {
		t.Fatalf("Tree has not size 2")
	}
	if tree.Height() != 2 {
		t.Fatalf("Tree has not height 2")
	}
}

func TestInsertContainsRight(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(5)
	tree.Insert(7)
	if !tree.Contains(7) {
		t.Fatalf("Tree does not contain 7")
	}
	if tree.Size() != 2 {
		t.Fatalf("Tree has not size 2")
	}
	if tree.Height() != 2 {
		t.Fatalf("Tree has not height 2")
	}
}

func TestInsertManyContainsAll(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{10, 3, 15, 1, 2, 100, 4, 17, 16, 9, 75, 8, 11, 12, 33, 20, 5, 6, 7, 22, 13, 14, 88, 18, 19}
	for i := range numbers {
		tree.Insert(numbers[i])
	}
	for i := range numbers {
		if !tree.Contains(numbers[i]) {
			t.Fatalf("Tree does not contain %d", i)
		}
	}
	if tree.Size() != len(numbers) {
		t.Fatalf("Tree does not have size %d", len(numbers))
	}
	if tree.Height() != 6 {
		t.Fatalf("Expected height 6, but got %d", tree.Height())
	}
	validateTreeProperties(t, tree.node)
	validateParentRefs(t, tree.node)
}

func TestStringEmpty(t *testing.T) {
	tree := MakeTree[int]()
	if tree.String() != "EmptyTree" {
		t.Fatalf("Empty Tree is not 'EmptyTree' but '%s'", tree.String())
	}
}

func TestStringSmallTree(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(5)
	tree.Insert(1)
	tree.Insert(7)
	expected := "L = {L = {} 1 R = {}} 5 R = {L = {} 7 R = {}}"
	if tree.String() != expected {
		t.Fatalf("Expected '%s' but is '%s'", expected, tree.String())
	}
}

func TestJsonEmpty(t *testing.T) {
	tree := MakeTree[int]()
	json := tree.Json()
	expected := "{}"
	if json != expected {
		t.Fatalf("Expected '%s', but is '%s'", expected, json)
	}
}

func TestJsonSmallTree(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(5)
	tree.Insert(1)
	tree.Insert(7)
	json := tree.Json()
	expected := `{"value": 5,"color": "black","parent": "nil","left": {"value": 1,"color": "red","parent": 5,"left": {},"right": {}},"right": {"value": 7,"color": "red","parent": 5,"left": {},"right": {}}}`
	if json != expected {
		t.Fatalf("Expected '%s', but is '%s'", expected, json)
	}
}

func TestRightRotatePanicWithoutLeftChild(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(5)
	tree.Insert(7)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Did not panic")
		}

	}()

	tree.node.RightRotate()
}

func TestRightRotateWithoutRightGrandchild(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{5, 3, 7, 1}
	for i := range numbers {
		tree.Insert(numbers[i])
	}

	tree.node.RightRotate()

	for i := range numbers {
		if !tree.Contains(numbers[i]) {
			t.Fatalf("Tree does not contain %d", i)
		}
	}
	if tree.Size() != len(numbers) {
		t.Fatalf("Tree does not have size %d", len(numbers))
	}
	validateParentRefs(t, tree.node)
}

func TestRightRotate(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{5, 3, 7, 4}
	for i := range numbers {
		tree.Insert(numbers[i])
	}

	tree.node.RightRotate()

	for i := range numbers {
		if !tree.Contains(numbers[i]) {
			t.Fatalf("Tree does not contain %d", i)
		}
	}
	if tree.Size() != len(numbers) {
		t.Fatalf("Tree does not have size %d", len(numbers))
	}
	validateParentRefs(t, tree.node)
}

func TestRightRotateNoUncle(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{1, 5, 15}
	for i := range numbers {
		tree.Insert(numbers[i])
	}

	tree.node.RightRotate()

	for i := range numbers {
		if !tree.Contains(numbers[i]) {
			t.Fatalf("Tree does not contain %d", i)
		}
	}
	if tree.Size() != len(numbers) {
		t.Fatalf("Tree does not have size %d", len(numbers))
	}
	validateParentRefs(t, tree.node)
}

func TestLeftRotateLine(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(1)
	tree.Insert(2)
	tree.Insert(3)

	if tree.node.val != 2 {
		t.Fatalf("2 did not become root")
	}
	if tree.node.right.val != 3 {
		t.Fatalf("3 did not become right child")
	}
	if tree.node.left.val != 1 {
		t.Fatalf("1 did not become left child")
	}
	if tree.node.right.p != tree.node {
		t.Fatalf("3's parent is not 2")
	}
	if tree.node.left.p != tree.node {
		t.Fatalf("1's parent is not 2")
	}
}

func TestLeftRotateNoUncle(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{15, 5, 1}
	for i := range numbers {
		tree.Insert(numbers[i])
	}

	tree.node.LeftRotate()

	for i := range numbers {
		if !tree.Contains(numbers[i]) {
			t.Fatalf("Tree does not contain %d", i)
		}
	}
	if tree.Size() != len(numbers) {
		t.Fatalf("Tree does not have size %d", len(numbers))
	}
	validateParentRefs(t, tree.node)
}

func TestLeftRotatePanicWithoutRightChild(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(7)
	tree.Insert(5)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Did not panic")
		}
	}()

	tree.node.LeftRotate()
}

func TestLeftRotateWithoutLeftGrandchild(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{5, 3, 7, 1}
	for i := range numbers {
		tree.Insert(numbers[i])
	}

	tree.node.LeftRotate()

	for i := range numbers {
		if !tree.Contains(numbers[i]) {
			t.Fatalf("Tree does not contain %d", i)
		}
	}
	if tree.Size() != len(numbers) {
		t.Fatalf("Tree does not have size %d", len(numbers))
	}
	validateParentRefs(t, tree.node)
}

func TestLeftRotate(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{5, 3, 7, 6}
	for i := range numbers {
		tree.Insert(numbers[i])
	}

	tree.node.LeftRotate()

	for i := range numbers {
		if !tree.Contains(numbers[i]) {
			t.Fatalf("Tree does not contain %d", i)
		}
	}
	if tree.Size() != len(numbers) {
		t.Fatalf("Tree does not have size %d", len(numbers))
	}
	validateParentRefs(t, tree.node)
}

func TestParentRefsEmpty(t *testing.T) {
	tree := MakeTree[int]()
	validateParentRefs(t, tree.node)
	validateTreeProperties(t, tree.node)
}

func TestParentRefsSmallTree(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{5, 3, 7, 1, 9, 2}

	for i := range numbers {
		tree.Insert(numbers[i])
		validateParentRefs(t, tree.node)
	}
}

func validateParentRefs[O cmp.Ordered](t *testing.T, n *Node[O]) {
	if n == nil {
		return
	}

	if n.p != nil && n.p.left != n && n.p.right != n {
		t.Fatalf("%s is not a child of their parent", String(n.val))
	}

	if n.left != nil && n.left.p != n {
		t.Fatalf("left child %s does not point back to parent", String(n.left.val))
	}

	if n.right != nil && n.right.p != n {
		t.Fatalf("right child %s does not point back to parent", String(n.left.val))
	}

	validateParentRefs(t, n.left)
	validateParentRefs(t, n.right)
}

func TestUncleLeftTriangle(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(4)
	tree.Insert(3)

	if uncle, rel, ok := tree.node.right.left.Uncle(); !ok || uncle.val != 1 || rel != Triangle {
		t.Fatalf("Uncle is not 1, or not Triangle")
	}
	validateTreeProperties(t, tree.node)
}

func TestUncleRightTriangle(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(3)
	tree.Insert(4)
	tree.Insert(1)
	tree.Insert(2)

	if uncle, rel, ok := tree.node.left.right.Uncle(); !ok || uncle.val != 4 || rel != Triangle {
		t.Fatalf("Uncle is not 4, or not Triangle")
	}
	validateTreeProperties(t, tree.node)
}

func TestUncleLeftLine(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(4)
	tree.Insert(5)

	if uncle, rel, ok := tree.node.right.right.Uncle(); !ok || uncle.val != 1 || rel != Line {
		t.Fatalf("Uncle is not 1, or not Line")
	}
	validateTreeProperties(t, tree.node)
}

func TestUncleRightLine(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(3)
	tree.Insert(4)
	tree.Insert(1)
	tree.Insert(0)

	if uncle, rel, ok := tree.node.left.left.Uncle(); !ok || uncle.val != 4 || rel != Line {
		t.Fatalf("Uncle is not 4, or not Line")
	}
	validateTreeProperties(t, tree.node)
}

func TestNoUncleRoot(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(1)

	if _, _, ok := tree.node.Uncle(); ok {
		t.Fatalf("Has an uncle")
	}
	validateTreeProperties(t, tree.node)
}

func TestNoUncleNoGrandparentLeft(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(2)
	tree.Insert(1)

	if _, _, ok := tree.node.left.Uncle(); ok {
		t.Fatalf("Has an uncle")
	}
	validateTreeProperties(t, tree.node)
}

func TestNoUncleNoGrandparentRight(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(1)
	tree.Insert(2)

	if _, _, ok := tree.node.right.Uncle(); ok {
		t.Fatalf("Has an uncle")
	}
	validateTreeProperties(t, tree.node)
}

func TestNoUncleLeft(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(1)

	validateTreeProperties(t, tree.node)
}

func TestNoUncleRight(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(1)
	tree.Insert(2)
	tree.Insert(3)

	validateTreeProperties(t, tree.node)
}

func TestNewInsertRootIsBlack(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(1)

	if tree.node.color != Black {
		t.Fatalf("root is not black")
	}

}

func TestRecolorUncleRed(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(3)
	tree.Insert(1)
	tree.Insert(4)

	if tree.node.right.color != Red {
		t.Fatalf("Precondition: right child should be red")

	}
	if tree.node.left.color != Red {
		t.Fatalf("Precondition: left child should be red")
	}

	tree.Insert(2)

	if tree.node.right.color != Black {
		t.Fatalf("Uncle was not recolored")
	}
	if tree.node.left.color != Black {
		t.Fatalf("Parent was not recolored")
	}
	if tree.node.color != Black {
		t.Fatalf("Grandparent was not recolored")
	}
	validateTreeProperties(t, tree.node)
}

func TestRotateLineLeft(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(4)
	tree.Insert(5)
	tree.Insert(2)

	tree.node.right.color = Black
	tree.node.left.color = Red

	tree.Insert(1)

	if tree.node.val != 2 {
		t.Fatalf("2 is not root")
	}
	if tree.node.color != Black {
		t.Fatalf("parent not recolored")
	}
	if tree.node.right.val != 4 {
		t.Fatalf("4 is not right child")
	}
	if tree.node.right.color != Red {
		t.Fatalf("Grandparent not recolored")
	}
	if tree.node.left.val != 1 {
		t.Fatalf("1 is not left child")
	}

}

func TestRotateLineRight(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(3)

	tree.node.left.color = Black
	tree.node.right.color = Red

	tree.Insert(4)

	if tree.node.val != 3 {
		t.Fatalf("3 is not root")
	}
	if tree.node.color != Black {
		t.Fatalf("Parent is not black")
	}
	if tree.node.left.val != 2 {
		t.Fatalf("2 is not left child")
	}
	if tree.node.left.color != Red {
		t.Fatalf("Grandparent is not red")
	}
	if tree.node.right.val != 4 {
		t.Fatalf("4 is not right child")
	}

}

func TestValidateTreePropertiesSmall(t *testing.T) {
	tree := MakeTree[int]()
	tree.Insert(2)
	tree.Insert(1)
	tree.Insert(3)

	validateTreeProperties(t, tree.node)
}

func TestValidateTreePropertiesIncreasing(t *testing.T) {
	tree := MakeTree[int]()
	for i := range 100 {
		tree.Insert(i)
	}
	validateTreeProperties(t, tree.node)
	validateParentRefs(t, tree.node)
}

func TestValidateTreePropertiesDecreasing(t *testing.T) {
	tree := MakeTree[int]()
	for i := range 100 {
		tree.Insert(100 - i)
	}
	validateTreeProperties(t, tree.node)
	validateParentRefs(t, tree.node)
}

func TestValidateTreePropertiesMedium(t *testing.T) {
	tree := MakeTree[int]()
	numbers := []int{10, 3, 15, 1, 2, 7, 11, 6}
	for i := range numbers {
		tree.Insert(numbers[i])
	}
	validateTreeProperties(t, tree.node)
}

func TestValidateTreePropertiesBig(t *testing.T) {
	tree := MakeTree[int]()
	for i := range 10000 {
		tree.Insert(i)
	}

	validateTreeProperties(t, tree.node)
}

func validateTreeProperties[O cmp.Ordered](t *testing.T, n *Node[O]) {
	if n == nil {
		return

	}

	if n.color != Black {
		t.Fatalf("Root is not black")
	}

	if err := validateRedNodeHasBlackChildren(n); err != nil {
		t.Fatal(err)
	}

	if _, err := validateSameNumberOfBlackNodesToLeaves(n); err != nil {
		t.Fatal(err)
	}

}

func validateRedNodeHasBlackChildren[O cmp.Ordered](n *Node[O]) error {
	if n == nil {
		return nil
	}

	if n.color == Red {
		if n.left != nil && n.left.color != Black {
			return fmt.Errorf("Red node %s has red left child %s", String(n.val), String(n.left.val))
		}
		if n.right != nil && n.right.color != Black {
			return fmt.Errorf("Red node %s has red right child %s", String(n.val), String(n.right.val))
		}
	}

	if err := validateRedNodeHasBlackChildren(n.left); err != nil {
		return err
	}

	if err := validateRedNodeHasBlackChildren(n.right); err != nil {
		return err
	}

	return nil
}

func validateSameNumberOfBlackNodesToLeaves[O cmp.Ordered](n *Node[O]) (int, error) {
	if n == nil {
		return 0, nil
	}

	left, err := validateSameNumberOfBlackNodesToLeaves(n.left)
	if err != nil {
		return 0, err
	}

	right, err := validateSameNumberOfBlackNodesToLeaves(n.right)
	if err != nil {
		return 0, err
	}

	if left != right {
		return 0, fmt.Errorf("Left child has %d black nodes, but right child has %d", left, right)
	}

	self := 0
	if n.color == Black {
		self = 1
	}

	return left + self, nil

}
