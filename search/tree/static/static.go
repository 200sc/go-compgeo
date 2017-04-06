package static

// Static is separated from the rest of tree
// because otherwise most of its types and functions
// would be preceded with "static"

type node struct {
	// eventually key should be a comparable interface
	// but that would probably poorly effect performance
	key float64
	val interface{}
}

// This implicitly says that
// a user cannot store nils in
// this tree. This is probably
// overly limiting.
func (n node) isNil() bool {
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

type BST []node

func (b BST) isNil(i int) bool {
	if i < len(b) {
		return b[i].isNil()
	}
	return true
}

func (b BST) minKey(i int) int {
	for !b.isNil(left(i)) {
		i = left(i)
	}
	return i
}

// Size :
// Return the number of elements in this static tree
func (b BST) Size() int {
	sz := 0
	b.size(0, &sz)
	return sz
}

// This is a goofy way of doing this.
// A better way would be to just keep track of size as we insert
// and delete, but this way we use less memory I suppose.
func (b BST) size(i int, sz *int) {
	if !b.isNil(i) {
		*sz++
		b.size(left(i), sz)
		b.size(right(i), sz)
	}
}

// Search :
// Search returns the first value of the given key found in the
// tree. No guarantee is made about what is returned if multiple
// nodes share the input key.
// If no key is found, returns (false, nil).
// Todo: benchmark value versus pointer receiver
func (b BST) Search(key float64) (bool, interface{}) {
	i := 1
	for {
		n := b[i]
		k := n.key
		if k == key {
			return true, n.val
		} else if k < key {
			i = left(i)
		} else {
			i = right(i)
		}
		if b.isNil(i) {
			return false, nil
		}
	}
}

// InOrderTraverse :
// There are multiple ways to traverse a tree.
// The most useful of these is the in-order traverse,
// and that's what we provide here.
// Other traversal methods can be added as needed.
func (b BST) InOrderTraverse() []node {
	out := make([]node, b.Size())
	i := 0
	b.inOrderTraverse(out, 0, &i)
	return out
}

func (b BST) inOrderTraverse(out []node, i int, nextIndex *int) {
	if !b.isNil(i) {
		b.inOrderTraverse(out, left(i), nextIndex)
		out[*nextIndex] = b[i]
		*nextIndex++
		b.inOrderTraverse(out, right(i), nextIndex)
	}
}
