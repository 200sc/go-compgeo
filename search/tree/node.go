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

func (n *node) deleteSwap(n2 *node) {
	// Delete the lower node's parent reference
	if n2.parent.left == n2 {
		n2.parent.left = nil
	} else {
		n2.parent.right = nil
	}
	// Point deleted node's children to the lifted node,
	// and vice versa.
	n2.left = n.left
	n2.left.parent = n2
	n2.right = n.right
	n2.right.parent = n2

	n2.parent = n.parent
	// Update deleted node's parent
	if n2.parent.left == n {
		n2.parent.left = n2
	} else {
		n2.parent.right = n2
	}
}
