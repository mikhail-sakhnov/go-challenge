package tree

type node struct {
	value  int
	height int
	left   *node
	right  *node
}

func newNode(value int) *node {
	return &node{
		value:  value,
		height: 1,
		left:   nil,
		right:  nil,
	}

}

func insert(root *node, value int) *node {
	if root == nil {
		return newNode(value)
	}
	// No equals branch prevents from doubles
	if value < root.value {
		root.left = insert(root.left, value)
	} else if value > root.value {
		root.right = insert(root.right, value)
	}
	return balance(root)
}

func toSlice(root *node) []int {
	result := []int{}
	stack := []*node{}
	current := root
	for {
		if current != nil {
			stack = append(stack, current)
			current = current.left
		} else {
			if len(stack) > 0 {
				current = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				result = append(result, current.value)
				current = current.right
			} else {
				break
			}
		}
	}
	return result
}

func height(n *node) int {
	if n == nil {
		return 0
	}
	return n.height
}

func factor(node *node) int {
	return height(node.right) - height(node.left)
}

func fixHeight(node *node) {
	if height(node.right) > height(node.left) {
		node.height = height(node.right) + 1
		return
	}
	node.height = height(node.left) + 1
}

func rotateRight(node *node) *node {
	q := node.left
	node.left = q.right
	q.right = node
	fixHeight(node)
	fixHeight(q)
	return q
}

func rotateLeft(node *node) *node {
	q := node.right
	node.right = q.left
	q.left = node
	fixHeight(q)
	fixHeight(node)
	return q
}

func balance(root *node) *node {
	fixHeight(root)
	if factor(root) == 2 {
		if factor(root.right) < 0 {
			root.right = rotateRight(root.right)
		}
		return rotateLeft(root)
	}
	if factor(root) == -2 {
		if factor(root.left) > 0 {
			root.left = rotateLeft(root.left)
		}
		return rotateRight(root)
	}
	return root
}
