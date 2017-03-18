package compgeo

type AVLNode struct {
	key float64
	val interface{}
}

type AVLTree []Node

// NewAVLTree :
func NewAVLTree() *AVLTree {
	t := new(AVLTree)
	// We don't have reason to do anything here
	return t
}

// Size :
func (avlt *AVLTree) Size() int {
	return len(*avlt)
}

// Insert :
func (avlt *AVLTree) Insert(n Node) error {
	return nil
}

// Delete :
// Because we allow duplicate keys,
// because real data has duplicate keys,
// we require you specify what you want to delete
// at the given key or nil if you know for sure that
// there is only one value with the given key.
func (avlt *AVLTree) Delete(n Node) error {
	return nil
}

// Search :
func (avlt *AVLTree) Search(key float64) (bool, interface{}) {
	t := *avlt
	i := 0
	for {
		n := t[i]
		k := n.Key()
		if k == key {
			return true, n.Val()
		} else if k < key {
			i = 2 * i
		} else {
			i = (2 * i) + 1
		}
		if i >= len(t) {
			return false, nil
		}
	}
}

// Traverse :
// There are multiple ways to traverse a tree.
// The most useful of these is the in-order traverse,
// and that's what we provide here.
// Other traversal methods can be added as needed.
func (avlt *AVLTree) Traverse() []Node {
	out := make([]Node, len(*avlt))
	out_i := 0
	avlt.inOrderTraverse(out, 0, &out_i)
	return out
}

func (avlt *AVLTree) inOrderTraverse(out []Node, i int, next_index *int) {
	if i < len(*avlt) {
		v := (*avlt)[i]
		// If a node is nil, it cannot have children
		if v != nil {
			avlt.inOrderTraverse(out, 2*i, next_index)
			out[*next_index] = v
			*next_index++
			avlt.inOrderTraverse(out, (2*i)+1, next_index)
		}
	}
}
