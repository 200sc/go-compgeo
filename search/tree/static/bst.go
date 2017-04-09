package static

import "github.com/200sc/go-compgeo/search"

// Static is separated from the rest of tree
// because otherwise most of its types and functions
// would be preceded with "static"

// BST in static is an array-based BST, as opposed
// to a pointer-based BST.
type BST []*Node

func (b *BST) isNil(i int) bool {
	bst := *b
	if i < len(bst) {
		return bst[i] == nil
	}
	return true
}

// could pseudo-binary search?
func (b *BST) minKey(i int) int {
	for !b.isNil(Left(i)) {
		i = Left(i)
	}
	return i
}

func (b *BST) maxKey(i int) int {
	for !b.isNil(Right(i)) {
		i = Right(i)
	}
	return i
}

// Size :
// Return the number of elements in this static tree
func (b *BST) Size() int {
	sz := 0
	b.size(0, &sz)
	return sz
}

// This is a goofy way of doing this.
// A better way would be to just keep track of size as we insert
// and delete, but this way we use less memory I suppose.
func (b *BST) size(i int, sz *int) {
	if !b.isNil(i) {
		*sz++
		b.size(Left(i), sz)
		b.size(Right(i), sz)
	}
}

func (b *BST) successor(i int) int {
	for b.isNil(i) {
		i = Parent(i)
	}
	if !b.isNil(Right(i)) {
		return b.minKey(Right(i))
	}
	j := Parent(i)
	for !b.isNil(j) && !isLeftChild(i) {
		i = j
		j = Parent(j)
	}
	return j
}

func (b *BST) predecessor(i int) int {
	for b.isNil(i) {
		i = Parent(i)
	}
	if !b.isNil(Left(i)) {
		return b.maxKey(Left(i))
	}
	j := Parent(i)
	for !b.isNil(j) && isLeftChild(i) {
		i = j
		j = Parent(j)
	}
	return j
}

// Search :
// Search returns the first value of the given key found in the
// tree. No guarantee is made about what is returned if multiple
// nodes share the input key.
// If no key is found, returns (false, nil).
//
// Value vs pointer reciever was benchmarked. Result: maybe pointer is better
func (b *BST) Search(key float64) (bool, interface{}) {
	i, ok := b.search(key)
	if ok {
		return true, (*b)[i].val
	}
	return false, nil
}

func (b *BST) SearchUp(key float64) interface{} {
	i, ok := b.search(key)
	if ok {
		return (*b)[i]
	}
	j := b.successor(i)
	bst := *b
	if b.isNil(j) || ((bst[j].key > bst[i].key) && (bst[i].key > key)) {
		j = i
	}
	return bst[j].val
}

func (b *BST) SearchDown(key float64) interface{} {
	i, ok := b.search(key)
	if ok {
		return (*b)[i]
	}
	j := b.predecessor(i)
	bst := *b
	if b.isNil(j) || ((bst[j].key < bst[i].key) && (bst[i].key < key)) {
		j = i
	}
	return bst[j].val
}

func (b *BST) search(key float64) (int, bool) {
	i := 1
	bst := *b
	var n *Node
	var k float64
	for {
		n = bst[i]
		k = n.key
		if k == key {
			return i, true
		}
		i = Left(i)
		if k < key {
			i++
		}
		if b.isNil(i) {
			return Parent(i), false
		}
	}
}

// InOrderTraverse :
// There are multiple ways to traverse a tree.
// The most useful of these is the in-order traverse,
// and that's what we provide here.
// Other traversal methods can be added as needed.
func (b *BST) InOrderTraverse() []search.Node {
	out := make([]search.Node, b.Size())
	i := 0
	b.inOrderTraverse(out, 0, &i)
	return out
}

func (b *BST) inOrderTraverse(out []search.Node, i int, nextIndex *int) {
	bst := *b
	if !b.isNil(i) {
		b.inOrderTraverse(out, Right(i), nextIndex)
		out[*nextIndex] = bst[i]
		*nextIndex++
		b.inOrderTraverse(out, Left(i), nextIndex)
	}
}
