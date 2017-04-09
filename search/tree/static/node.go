package static

// Node is a data structures which
// just contain a key and value
type Node struct {
	// eventually key should be a comparable interface
	// but that would probably poorly effect performance
	key float64
	val interface{}
}

// *Node vs Node was benchmarked.
// Result: maybe * is better. It's hard to tell.

// NewNode returns a node constructed from the input
// key and value.
func NewNode(k float64, v interface{}) *Node {
	return &Node{key: k, val: v}
}

// Key returns the key of this node.
func (n Node) Key() float64 {
	return n.key
}

// Val returns the value of this node.
func (n Node) Val() interface{} {
	return n.val
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
