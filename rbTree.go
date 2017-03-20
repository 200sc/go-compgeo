package compgeo

func NewRBTree() *BST {
	t := NewBST()
	t.typ = RBTreeType
	return t
}

func color(n Node) bool {
	return n.payload.(bool)
}

func (bst *BST) color(i int) bool {
	return color(bst.tree[i])
}

func (bst *BST) rb_delete(i int) {
	// If this node has two children
	if !bst.IsNil(right(i)) && !bst.IsNil(left(i)) {
		// Find the maximum value of the left subtree
		// or the minimum value of the right subtree.
		// Presumably defaultin to one over the other will
		// cause the tree to lean in one direction over the
		// other. Needs to be tested.

		// if rand.Float64() < 0.5 {
		j := bst.minKey(right(i))
		//} else {
		// j := bst.maxKey(left(i))
		//}
		bst.tree[i], bst.tree[j] = bst.tree[j], bst.tree[i]
		i = j
	}
}

func (bst *BST) rb_insert(i int) {
	// If i is the root
	for {
		p := parent(i)
		if i == 1 {
			bst.tree[1].payload = BLACK
			return
		}
		// i's parent must exist, as i is not the root ---
		// If i's parent is black
		if bst.color(p) == BLACK {
			return
		}

		// i's grandparent must exist, as i's parent is red. ---
		// if i's grandparent did not exist, i's parent would
		// be the root and would be black.
		// If i's parent is red and i's uncle is red

		redUncle := false
		gp := ancestor(i, 2)
		// i's parent is a left child if it is even
		if isLeftChild(p) {
			if bst.color(right(gp)) == RED {
				redUncle = true
			}
		} else if bst.color(left(gp)) == RED {
			redUncle = true
		}
		// Would this be faster?
		// if bst.color(right(gp)) == RED && bst.color(left(gp)) == RED {...}
		if redUncle {
			i = gp
			bst.tree[left(i)].payload = BLACK
			bst.tree[right(i)].payload = BLACK
			bst.tree[i].payload = RED
			// recurse on i's grandparent

		} else {
			bst.tree[gp].payload = RED
			if left(gp) == p {
				// if i is a right child
				if !isLeftChild(i) {
					bst.leftRotate(p)
				}
				bst.tree[p].payload = BLACK
				bst.rightRotate(gp)
			} else {
				if isLeftChild(i) {
					bst.rightRotate(p)
				}
				bst.tree[p].payload = BLACK
				bst.leftRotate(gp)
			}
			return
		}
	}
}
