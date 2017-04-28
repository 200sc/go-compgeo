package static

import "github.com/200sc/go-compgeo/search"

// Node is a data structures which
// just contain a key and value
type Node struct {
	key search.Comparable
	// Unlike pointer BSTs, right now static BSTs don't support
	// multiple-valued keys, as our API only has a difference in
	// how they are dealt with in modification cases. In search
	// cases, as search just takes a key (right now), we always
	// return the same value.
	val search.Equalable
}

// *Node vs Node was benchmarked.
// Result: maybe * is better. It's hard to tell.

// NewNode returns a node constructed from the input
// key and value.
func NewNode(k search.Comparable, v search.Equalable) *Node {
	return &Node{key: k, val: v}
}

// Key returns the key of this node.
func (n Node) Key() search.Comparable {
	return n.key
}

// Val returns the value of this node.
func (n Node) Val() search.Equalable {
	return n.val
}

func (n Node) copy() *Node {
	return &Node{n.key, n.val}
}

// This implicitly says that
// a user cannot store nils in
// this tree. This is probably
// overly limiting.
func (n Node) isNil() bool {
	return n.val == nil
}

// Ancestor returns the index representing
// this node's ancestor of N generations
// equivalent to calling Parent tiersUp times
func Ancestor(i, tiersUp int) int {
	return i / (tiersUp * 2)
}

// Parent returns the index represeting
// this node's direct parent.
func Parent(i int) int {
	return i / 2
}

// Left returns the index representing this
// node's left child.
func Left(i int) int {
	return 2 * i
}

// Right returns the index representing this node's
// right child.
func Right(i int) int {
	return (2 * i) + 1
}

// This is an interesting, useful thing that would be
// useful if the static tree cared.
func isLeftChild(i int) bool {
	return i%2 == 0
}
