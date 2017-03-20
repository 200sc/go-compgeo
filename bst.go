package compgeo

import "errors"

// BST Type enumerator
const (
	AVLTreeType = iota
	RBTreeType
	TTreeType
	SplayTreeType
	TangoTreeType
)

const (
	RED   = false
	BLACK = true
)

type Node struct {
	// eventually key should be a comparable interface
	// but that would probably poorly effect performance
	key float64
	val interface{}
	// Each tree type has a different payload (if any)
	payload interface{}
}

// This implicitly says that
// a user cannot store nils in
// this tree. This is probably
// overly limiting.
func (n Node) IsNil() bool {
	return n.val == nil
}

func ancestor(i, tiersUp int) int {
	return i / (tiersUp * 2)
}

func parent(i int) int {
	return i / 2
}

func left(i int) int {
	return 2 * i
}

func right(i int) int {
	return (2 * i) + 1
}

func isLeftChild(i int) bool {
	return i%2 == 0
}

type BST struct {
	tree []Node
	typ  int
	size int
}

// NewBST :
func NewBST() *BST {
	bst := new(BST)
	// The 0th element of this tree represents
	// going off the top of the tree.
	bst.tree = append(bst.tree, Node{})
	// We don't have reason to do anything here
	return bst
}

func (bst *BST) IsNil(i int) bool {
	if i < len(bst.tree) {
		return bst.tree[i].IsNil()
	}
	return true
}

// Size :
func (bst *BST) Size() int {
	return bst.size
}

// Insert :
func (bst *BST) Insert(n Node) error {
	t := bst.tree
	i := 1
	for {
		if i >= len(t) || n.IsNil() {
			break
		}
		if t[i].key < n.key {
			i = left(i)
		} else {
			i = right(i)
		}
	}
	if i >= len(t) {
		bst.tree = append(bst.tree, (make([]Node, 1+(i-len(t))))...)
	}
	// Lets say this is a RB tree

	bst.tree[i] = n
	bst.size++

	switch bst.typ {
	case RBTreeType:
		bst.tree[i].payload = RED
		bst.rb_insert(i)
	}
	return nil
}

// Delete :
// Because we allow duplicate keys,
// because real data has duplicate keys,
// we require you specify what you want to delete
// at the given key or nil if you know for sure that
// there is only one value with the given key (or
// do not care what is deleted).
func (bst *BST) Delete(n Node) error {
	t := bst.tree
	i := 1
	for {
		n2 := t[i]
		k := n2.key
		if k == n.key {
			if n.val == nil {
				break
			}
			for t[i].val != n.val {
				i = right(i)
				if i >= len(t) {
					return errors.New("Value not found")
				}
			}
			break
		} else if k < n.key {
			i = left(i)
		} else {
			i = right(i)
		}
		if i >= len(t) {
			return errors.New("Key not found")
		}
	}
	switch bst.typ {
	case RBTreeType:
		bst.rb_delete(i)
	}
	return nil
}

// Search :
func (bst *BST) Search(key float64) (bool, interface{}) {
	t := bst.tree
	i := 1
	for {
		n := t[i]
		k := n.key
		if k == key {
			return true, n.val
		} else if k < key {
			i = left(i)
		} else {
			i = right(i)
		}
		if i >= len(t) {
			return false, nil
		}
	}
}

func (bst *BST) minKey(i int) int {
	for !bst.IsNil(left(i)) {
		i = left(i)
	}
	return i
}

// Traverse :
// There are multiple ways to traverse a tree.
// The most useful of these is the in-order traverse,
// and that's what we provide here.
// Other traversal methods can be added as needed.
func (bst *BST) InOrderTraverse() []Node {
	out := make([]Node, len(bst.tree))
	out_i := 0
	bst.inOrderTraverse(out, 0, &out_i)
	return out
}

func (bst *BST) inOrderTraverse(out []Node, i int, next_index *int) {
	if i < len(bst.tree) {
		v := bst.tree[i]
		// If a node is nil, it cannot have children
		if !v.IsNil() {
			bst.inOrderTraverse(out, left(i), next_index)
			out[*next_index] = v
			*next_index++
			bst.inOrderTraverse(out, right(i), next_index)
		}
	}
}
