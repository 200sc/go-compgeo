package tree

import "errors"

const (
	red   = false
	black = true
)

var (
	rbFnSet = &fnSet{
		insertFn: rbInsert,
		deleteFn: rbDelete,
		searchFn: rbSearch,
	}
)

func (n *node) isBlack() bool {
	if n == nil {
		return true
	}
	return (n.payload.(bool) == black)
}

func (n *node) ancestor(i int) *node {
	for j := 0; j < i; j++ {
		if n == nil {
			return n
		}
		n = n.parent
	}
	return n
}

// RBValid returns whether the given BST is a valid Red Black tree
func RBValid(bst *BST) (bool, error) {
	n := bst.root
	if n == nil {
		return true, nil
	}
	// We satisfy case 3, that the leaves must be black,
	// implicitly as we evalaute nil to be black.
	// Case 2: the root must be black
	switch n.payload.(type) {
	case bool:
		if n.payload == red {
			return false, errors.New("The root is not black")
		}
	}
	b, _, err := n.RBValid(true)
	return b, err
}

// RBValid returns whether the given node is a valid Red Black Subtree.
// It returns boolean validity, a potential error (if b = false, err = nil)
// and the number of black nodes on any path starting from it.
func (n *node) RBValid(mustBeBlack bool) (bool, int, error) {
	if n != nil {
		switch n.payload.(type) {
		case bool:
			mustBeBlack = false
			increaseCt := 0
			if n.payload == red {
				if mustBeBlack {
					return false, 0, errors.New("A red node's child was red")
				}
				mustBeBlack = true
			} else {
				increaseCt = 1
			}
			b, ct1, err := n.left.RBValid(mustBeBlack)
			if !b {
				return b, 0, err
			}
			b, ct2, err := n.right.RBValid(mustBeBlack)
			if !b {
				return b, 0, err
			}
			if ct1 != ct2 {
				return false, 0, errors.New("The count of black nodes at either side of a subtree was not the same")
			}
			return true, ct1 + increaseCt, nil
		// Case 1: Each node is red or black
		default:
			return false, 0, errors.New("A node was neither red nor black")
		}
	}
	return true, 1, nil
}

func rbSearch(n *node) {
	// NOP
}

func rbInsert(n *node) {
	// If i is the root
	for {
		p := n.parent
		if p == nil {
			n.payload = black
			return
		}
		// i's parent must exist, as i is not the root ---
		// If i's parent is black
		if p.isBlack() {
			return
		}

		// i's grandparent must exist, as i's parent is red. ---
		// if i's grandparent did not exist, i's parent would
		// be the root and would be black.
		// If i's parent is red and i's uncle is red

		gp := p.parent
		uncle := n.uncle()
		if !uncle.isBlack() {
			gp.left.payload = black
			gp.right.payload = black
			gp.payload = red
			n = gp
			// Recurse
		} else {
			if p.right == n && p == gp.left {
				p.leftRotate()
				n = n.left
			} else if p.left == n && p == gp.right {
				p.rightRotate()
				n = n.right
			}
			p = n.parent
			gp = p.parent

			p.payload = black
			gp.payload = red

			if p.left == n {
				gp.rightRotate()
			} else {
				gp.leftRotate()
			}
			return
		}
	}
}
func rbDelete(n *node) {

	c := n.payload
	var r *node
	var newRoot *node
	if n.right == nil {
		r = n.left
		n.parentReplace(n.left)
	} else if n.left == nil {
		r = n.right
		n.parentReplace(n.right)
	} else {
		// Find the maximum value of the left subtree
		// or the minimum value of the right subtree.
		// Presumably defaulting to one over the other will
		// cause the tree to lean in one direction over the
		// other.

		// if rand.Float64() < 0.5 {
		n2 := n.right.minKey()
		//} else {
		// n2 := n.left.maxKey()
		//}
		r := n2.right
		if n2.parent == n {
			if r != nil {
				r.parent = n2
			}
		} else {
			n2.parentReplace(r)
			n2.right = n.right
			n2.right.parent = n2
		}
		if n.parent == nil {
			newRoot = n2
		}
		n.parentReplace(n2)
		if newRoot == n2 {
			n2.parent = nil
		}
		n2.left = n.left
		n2.left.parent = n2
		n2.payload = n.payload
	}
	if r == nil {
		return
	}
	if c.(bool) == black {
		n = r
		p := n.parent
		var s *node
		for p != nil && n.isBlack() {
			if n == p.left {
				s = p.right
				if !s.isBlack() {
					p.payload = red
					s.payload = black
					p.leftRotate()
					s = parent_sibling(n, p)
				}
				if s != nil {
					if s.right.isBlack() {
						if s.left.isBlack() {
							s.payload = red
							p = p.parent
							n = p
						} else {
							s.left.payload = black
							s.payload = red
							s.leftRotate()
							s = parent_sibling(n, p)
						}
					}
					if !s.right.isBlack() {
						s.payload = p.payload
						p.payload = black
						s.right.payload = black
						p.leftRotate()
						return
					}
				}
			} else {
				s = p.left
				if !s.isBlack() {
					p.payload = red
					s.payload = black
					p.rightRotate()
					s = parent_sibling(n, p)
				}
				if s != nil {
					if s.left.isBlack() {
						if s.right.isBlack() {
							s.payload = red
							p = p.parent
							n = p
						} else {
							s.right.payload = black
							s.payload = red
							s.leftRotate()
							s = parent_sibling(n, p)
						}
					}
					if !s.left.isBlack() {
						s.payload = p.payload
						p.payload = black
						s.left.payload = black
						p.rightRotate()
						return
					}
				}
			}
		}
	}
	if n != nil {
		n.payload = black
	}
}
