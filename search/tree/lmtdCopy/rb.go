package lmtdCopy

// import "errors"

// const (
// 	red   = false
// 	black = true
// )

// var (
// 	// RbFnSet performs RB insert and
// 	// RB delete for inserts and deletes,
// 	// and does nothing on lookups.
// 	RbFnSet = &FnSet{
// 		InsertFn: rbInsert,
// 		DeleteFn: rbDelete,
// 		SearchFn: nopNode,
// 	}
// )

// // For readability
// func (n *node) isRed() bool {
// 	return !n.isBlack()
// }

// func (n *node) isBlack() bool {
// 	if n == nil {
// 		return true
// 	}
// 	return n.payload.(bool) == black
// }

// // RBValid returns whether the given BST is a valid Red Black tree
// func RBValid(bst *BST) (bool, error) {
// 	n := bst.root
// 	if n == nil {
// 		return true, nil
// 	}
// 	// We satisfy case 3, that the leaves must be black,
// 	// implicitly as we evalaute nil to be black.
// 	// Case 2: the root must be black
// 	switch n.payload.(type) {
// 	case bool:
// 		if n.payload == red {
// 			return false, errors.New("The root is not black")
// 		}
// 	}
// 	b, _, err := n.RBValid(true)
// 	return b, err
// }

// // RBValid returns whether the given node is a valid Red Black Subtree.
// // It returns boolean validity, a potential error (if b = false, err = nil)
// // and the number of black nodes on any path starting from it.
// func (n *node) RBValid(mustBeBlack bool) (bool, int, error) {
// 	if n != nil {
// 		switch n.payload.(type) {
// 		case bool:
// 			mustBeBlack = false
// 			increaseCt := 0
// 			if n.payload == red {
// 				if mustBeBlack {
// 					return false, 0, errors.New("A red node's child was red")
// 				}
// 				mustBeBlack = true
// 			} else {
// 				increaseCt = 1
// 			}
// 			b, ct1, err := n.left.RBValid(mustBeBlack)
// 			if !b {
// 				return b, 0, err
// 			}
// 			b, ct2, err := n.right.RBValid(mustBeBlack)
// 			if !b {
// 				return b, 0, err
// 			}
// 			if ct1 != ct2 {
// 				return false, 0, errors.New("The count of black nodes at either side of a subtree was not the same")
// 			}
// 			return true, ct1 + increaseCt, nil
// 		// Case 1: Each node is red or black
// 		default:
// 			return false, 0, errors.New("A node was neither red nor black")
// 		}
// 	}
// 	return true, 1, nil
// }

// func rbInsert(n *node) (newRoot *node) {
// 	for {
// 		p := n.parent
// 		if p == nil {
// 			n.payload = black
// 			return
// 		}
// 		// i's parent must exist, as i is not the root ---
// 		// If i's parent is black
// 		if p.isBlack() {
// 			return
// 		}

// 		// i's grandparent must exist, as i's parent is red. ---
// 		// if i's grandparent did not exist, i's parent would
// 		// be the root and would be black.
// 		// If i's parent is red and i's uncle is red

// 		gp := p.parent
// 		uncle := n.uncle()
// 		if !uncle.isBlack() {
// 			gp.left.payload = black
// 			gp.right.payload = black
// 			gp.payload = red
// 			n = gp
// 			// Recurse
// 		} else {
// 			if p.right == n && p == gp.left {
// 				newRoot = root(newRoot, p.leftRotate())
// 				n = n.left
// 			} else if p.left == n && p == gp.right {
// 				newRoot = root(newRoot, p.rightRotate())

// 				n = n.right
// 			}
// 			p = n.parent
// 			gp = p.parent

// 			p.payload = black
// 			gp.payload = red

// 			if p.left == n {
// 				newRoot = root(newRoot, gp.rightRotate())
// 			} else {
// 				newRoot = root(newRoot, gp.leftRotate())
// 			}
// 			return
// 		}
// 	}
// }

// func rbDelete(n *node) (newRoot *node) {

// 	var c bool
// 	c = n.payload.(bool)
// 	var r *node
// 	//var newRoot *node
// 	p := n.parent
// 	if n.right == nil {
// 		r = n.left
// 		newRoot = n.parentReplace(n.left)
// 	} else if n.left == nil {
// 		r = n.right
// 		newRoot = n.parentReplace(n.right)
// 	} else {
// 		// Find the maximum value of the left subtree
// 		// or the minimum value of the right subtree.
// 		// Presumably defaulting to one over the other will
// 		// cause the tree to lean in one direction over the
// 		// other.

// 		// if rand.Float64() < 0.5 {
// 		n2 := n.right.minKey()
// 		c = n2.payload.(bool)
// 		//} else {
// 		// n2 := n.left.maxKey()
// 		//}
// 		p = n2.parent
// 		r = n2.right
// 		if n2.parent == n {
// 			if r != nil {
// 				r.parent = n2
// 			} else {
// 				p = n2
// 			}
// 		} else {
// 			newRoot = n2.parentReplace(r)
// 			n2.right = n.right
// 			n2.right.parent = n2
// 		}
// 		newRoot = root(newRoot, n.parentReplace(n2))
// 		n2.left = n.left
// 		n2.left.parent = n2
// 		n2.payload = n.payload
// 		if p == n {
// 			p = n2
// 		}
// 	}
// 	if c == black {
// 		newRoot = root(newRoot, rbDeleteFixup(r, p))
// 	}
// 	return
// }

// // DeleteFixup takes n and p, as nil nodes do
// // not contain a reference to their parent.
// // Note on cyclomtic complexity: RB delete fixup
// // cases don't have intuitive names, they're generally
// // referred to as case_N or fixup_N for n = 1..6.
// // Instead of making a bunch of numbered functions,
// // this implementation prefers to keep everything together
// // (as is common).
// func rbDeleteFixup(n, p *node) (newRoot *node) {
// 	var s *node
// 	for n.isBlack() {
// 		if n != nil {
// 			p = n.parent
// 		}
// 		// Case 1: p = nil
// 		// n is the root.
// 		if p == nil {
// 			newRoot = n
// 			break
// 		}
// 		// The subtree P->N has one fewer black nodes than P->S.
// 		s = pSibilng(n, p)
// 		if s.isRed() {
// 			// Case 2
// 			// S is red, so P is black.
// 			//
// 			// Give N a Black Sibling and
// 			// a Red parent.
// 			//
// 			p.payload = red
// 			s.payload = black
// 			if s == p.right {
// 				newRoot = root(p.leftRotate(), newRoot)
// 				s = p.right
// 			} else {
// 				newRoot = root(p.rightRotate(), newRoot)
// 				s = p.left
// 			}
// 			// Now P->N = P->NewS - 1, still,
// 			// and OldS->P = OldS-> p.sibling - 1
// 			//
// 			// R:P
// 			// |-- B:N
// 			// |-- B:S, not nil
// 		}
// 		//
// 		// Case 2.3: S is nil
// 		// We think this is impossible
// 		// if s == nil {
// 		// 	break
// 		// }
// 		// Case 3: Everything is black
// 		// In this case, Because S's children are black we can turn it red.
// 		// This means P->S = P->N, but GP->P = GP->P's sibling - 1,
// 		// so we recurse with n = p, p = gp.
// 		// --we crashed here once!!!?
// 		if p.isBlack() && s.isBlack() && s.left.isBlack() && s.right.isBlack() {
// 			s.payload = red
// 			n = p
// 			p = n.parent
// 			continue
// 		}
// 		// Case 4: Everything but P is black.
// 		// We can turn S red here as well, if we also make P red.
// 		// That will make P->N = P->S and they'll both be what they were
// 		// before the deletion, so we're done.
// 		if p.isRed() && s.isBlack() && s.left.isBlack() && s.right.isBlack() {
// 			s.payload = red
// 			p.payload = black
// 			break
// 		}
// 		// Case 5.1:
// 		// S has a left red child and a right black child,
// 		// and n is P's left child. A rotation will convert this
// 		// to case 6.
// 		if n == p.left && s.right.isBlack() && s.left.isRed() {
// 			s.payload = red
// 			s.left.payload = black
// 			newRoot = root(s.rightRotate(), newRoot)
// 			s = p.right
// 			// Case 5.2:
// 			// As 5.1, but flipped
// 		} else if n == p.right && s.left.isBlack() && s.right.isRed() {
// 			s.payload = red
// 			s.right.payload = black
// 			newRoot = root(s.leftRotate(), newRoot)
// 			s = p.left
// 		}
// 		// Case 6:
// 		// ...

// 		s.payload = p.payload
// 		p.payload = black
// 		if n == p.left {
// 			s.right.payload = black
// 			newRoot = root(p.leftRotate(), newRoot)
// 		} else {
// 			s.left.payload = black
// 			newRoot = root(p.rightRotate(), newRoot)
// 		}
// 		break
// 	}
// 	if n != nil {
// 		n.payload = black
// 	}
// 	return
// }

// func root(n1, n2 *node) *node {
// 	if n1 == nil {
// 		return n2
// 	}
// 	return n1
// }
