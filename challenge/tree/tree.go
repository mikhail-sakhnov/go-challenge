package tree

// Tree is wrapper over tree structure to have more oop like api
type Tree struct {
	root *node
}

// New builds new avl tree
func New() *Tree {
	return &Tree{}
}

// Insert inserts value into sorted tree
func (t *Tree) Insert(v int) {
	if t.root == nil {
		t.root = newNode(v)
		return
	}
	t.root = insert(t.root, v)
}

// ToSlice returns slice (dfs in-order path)
func (t *Tree) ToSlice() []int {
	return toSlice(t.root)
}
