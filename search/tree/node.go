package tree

type node struct {
	// eventually key should be a comparable interface
	// but that would probably poorly effect performance
	key float64
	val interface{}
	// Each tree type might have a different payload on each node
	// a good example of this is RED or BLACK on RBtrees.
	payload interface{}

	left, right, parent *node
}

func (n *node) Key() float64 {
	return n.key
}

func (n *node) Val() interface{} {
	return n.val
}

func (n *node) minKey() *node {
	if n.left == nil {
		return n
	}
	return n.left.minKey()
}

func (n *node) maxKey() *node {
	if n.right == nil {
		return n
	}
	return n.right.maxKey()
}

func (n *node) sibling() *node {
	if n.parent == nil {
		return nil
	}
	if n.parent.left == n {
		return n.parent.right
	}
	return n.parent.left
}

// Replace n.parent's pointer to n
// with a pointer to n2
func (n *node) parentReplace(n2 *node) {
	n2.parent = n.parent
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = n2
		} else {
			n.parent.right = n2
		}
	}
}

func (n *node) leftRotate() {
	newRight := n.right.left
	n.right.parent = n.parent
	if n.parent.left == n {
		n.parent.left = n.right
	} else {
		n.parent.right = n.right
	}
	n.right.left = n
	n.parent = n.right
	n.right = newRight
}

func (n *node) rightRotate() {
	// I would panic on n.left (or n)
	// being nil, but the panic will
	// happen anyway on trying to access
	// an element of nil.
	newLeft := n.left.right
	n.left.parent = n.parent
	if n.parent.left == n {
		n.parent.left = n.left
	} else {
		n.parent.right = n.left
	}
	n.left.right = n
	n.parent = n.left
	n.left = newLeft
}

func (n *node) swap(n2 *node) {
	// Reset the second node's parent reference
	if n2.parent != nil {
		if n2.parent.left == n2 {
			n2.parent.left = n
		} else {
			n2.parent.right = n
		}
	}
	n2parent := n2.parent
	n2left := n2.left
	n2right := n2.right

	// Point deleted node's children to the lifted node,
	// and vice versa.
	n2.left = n.left
	n2.right = n.right
	if n2.left != nil {
		n2.left.parent = n2
	}
	if n2.right != nil {
		n2.right.parent = n2
	}

	n2.parent = n.parent
	// Repeat for the first node
	if n2.parent != nil {
		if n2.parent.left == n {
			n2.parent.left = n2
		} else {
			n2.parent.right = n2
		}
	}

	n.left = n2left
	n.right = n2right
	if n.left != nil {
		n.left.parent = n
	}
	if n.right != nil {
		n.right.parent = n
	}

	n.parent = n2parent
}
