package static

type Node struct {
	// eventually key should be a comparable interface
	// but that would probably poorly effect performance
	key float64
	val interface{}
}

func NewNode(k float64, v interface{}) Node {
	return Node{key: k, val: v}
}

func (n Node) Key() float64 {
	return n.key
}

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

func Ancestor(i, tiersUp int) int {
	return i / (tiersUp * 2)
}

func Parent(i int) int {
	return i / 2
}

func Left(i int) int {
	return 2 * i
}

func Right(i int) int {
	return (2 * i) + 1
}

func isLeftChild(i int) bool {
	return i%2 == 0
}
