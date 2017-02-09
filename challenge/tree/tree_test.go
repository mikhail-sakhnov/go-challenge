package tree

import (
	"reflect"
	"testing"
)

// TestSmokeTestForTree smoke tests checks that we can get traverse from tree and the tree balances itself
func TestSmokeTestForTree(t *testing.T) {
	tree := New()
	tree.Insert(10)
	if factor(tree.root) > 2 {
		t.Fatalf("Height of %v is ", tree, tree.root.height)
	}
	tree.Insert(3)
	if factor(tree.root) > 2 {
		t.Fatalf("Height of %v is ", tree, tree.root.height)
	}
	tree.Insert(5)
	if factor(tree.root) > 2 {
		t.Fatalf("Height of %v is ", tree, tree.root.height)
	}
	tree.Insert(6)
	if factor(tree.root) > 2 {
		t.Fatalf("Height of %v is ", tree, tree.root.height)
	}
	tree.Insert(8)
	if factor(tree.root) > 2 {
		t.Fatalf("Height of %v is ", tree, tree.root.height)
	}
	tree.Insert(7)
	if factor(tree.root) > 2 {
		t.Fatalf("Height of %v is ", tree, tree.root.height)
	}
	tree.Insert(2)

	if factor(tree.root) > 2 {
		t.Fatalf("Height of %v is ", tree, tree.root.height)
	}
	tree.Insert(2) // We should have preventing doubles by design
	res := tree.ToSlice()
	exp := []int{2, 3, 5, 6, 7, 8, 10}
	if !reflect.DeepEqual(res, exp) {
		t.Fatalf("Expected %v, got %v", exp, res)
	}
}
