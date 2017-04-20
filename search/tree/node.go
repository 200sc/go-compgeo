package tree

import (
	"fmt"

	"github.com/200sc/go-compgeo/printutil"
	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree/static"
)

type node struct {
	// eventually key should be a comparable interface
	// but that would probably poorly effect performance
	key search.Comparable
	val []search.Equalable
	// Each tree type might have a different payload on each node
	// a good example of this is RED or BLACK on RBtrees.
	payload interface{}

	left, right, parent *node
}

func (n *node) calcSize() int {
	if n == nil {
		return 0
	}
	return n.left.calcSize() + n.right.calcSize() + len(n.val)
}

func (n *node) Key() search.Comparable {
	return n.key
}

func (n *node) Val() search.Equalable {
	return n.val[0]
}

func (n *node) isValid() (bool, search.Comparable, search.Comparable) {
	if n == nil {
		return true, search.NegativeInf{}, search.Inf{}
	}
	ok, min, max2 := n.left.isValid()
	if !ok {
		return false, nil, nil
	}
	ok, min2, max := n.right.isValid()
	if !ok {
		return false, nil, nil
	}

	if n.key.Compare(min) == search.Less ||
		n.key.Compare(max) == search.Greater {
		return false, nil, nil
	}
	if min2.Compare(min) == search.Less {
		min = min2
	}
	if max2.Compare(max) == search.Greater {
		max = max2
	}
	return true, min, max
}

func (n *node) copy() *node {
	if n == nil {
		return nil
	}
	cp := new(node)
	cp.left = n.left.copy()
	cp.right = n.right.copy()

	cp.key = n.key
	cp.val = make([]search.Equalable, len(n.val))
	copy(cp.val, n.val)

	if cp.left != nil {
		cp.left.parent = cp
	}
	if cp.right != nil {
		cp.right.parent = cp
	}
	cp.payload = n.payload

	return cp
}

func (n *node) minKey() *node {
	if n.left == nil {
		return n
	}
	return n.left.minKey()
}

func (n *node) maxKey() *node {
	if n.right == nil {
		return n
	}
	return n.right.maxKey()
}

func (n *node) successor() *node {
	if n == nil {
		return nil
	}
	if n.right != nil {
		return n.right.minKey()
	}
	p := n.parent
	for p != nil && n == p.right {
		n = p
		p = p.parent
	}
	return p
}

func (n *node) predecessor() *node {
	if n == nil {
		return nil
	}
	if n.left != nil {
		return n.left.maxKey()
	}
	p := n.parent
	for p != nil && n == p.left {
		n = p
		p = p.parent
	}
	return p
}

func (n *node) sibling() *node {
	p := n.parent
	if p == nil {
		return nil
	}
	if p.left == n {
		return p.right
	}
	return p.left
}

func pSibilng(n, p *node) *node {
	if p.left == n {
		return p.right
	}
	return p.left
}

func (n *node) uncle() *node {
	p := n.parent
	if p == nil {
		return nil
	}
	return p.sibling()
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

// Replace n.parent's pointer to n
// with a pointer to n2
func (n *node) parentReplace(n2 *node) *node {
	// if n.parent is nil, that means this is the root!
	// we're removing n from the tree, and our method of
	// finding then new root when a root is removed is to
	// follow the pointer of the old root. SO--
	var toReturn *node
	if n.parent == nil {
		toReturn = n2
	} else if n.parent.left == n {
		n.parent.left = n2
	} else {
		n.parent.right = n2
	}
	if n2 != nil {
		n2.parent = n.parent
	}
	return toReturn
}

func (n *node) leftRotate() (newRoot *node) {
	r := n.right
	n.right = r.left
	if r.left != nil {
		r.left.parent = n
	}
	r.parent = n.parent
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = r
		} else {
			n.parent.right = r
		}
	} else {
		newRoot = r
	}
	r.left = n
	n.parent = r
	return
}

func (n *node) rightRotate() (newRoot *node) {
	l := n.left
	n.left = l.right
	if l.right != nil {
		l.right.parent = n
	}
	l.parent = n.parent
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = l
		} else {
			n.parent.right = l
		}
	} else {
		newRoot = l
	}
	l.right = n
	n.parent = l
	return
}

func (n *node) staticTree(m map[int]*static.Node, i int) (map[int]*static.Node, int) {
	if n == nil {
		return m, 0
	}
	m[i] = static.NewNode(n.key, n.val[0])
	var maxIndex1, maxIndex2 int
	m, maxIndex1 = n.left.staticTree(m, static.Left(i))
	m, maxIndex2 = n.right.staticTree(m, static.Right(i))
	if maxIndex1 < maxIndex2 {
		maxIndex1 = maxIndex2
	}
	if maxIndex1 < i {
		maxIndex1 = i
	}
	return m, maxIndex1
}

func inOrderTraverse(n *node) []search.Node {
	if n != nil {
		lst := inOrderTraverse(n.left)
		lst = append(lst, n)
		return append(lst, inOrderTraverse(n.right)...)
	}
	return []search.Node{}
}

func (n *node) String() string {
	return n.string("", true)
}
func (n *node) string(prefix string, isTail bool) string {
	if n == nil || len(prefix) > 64 {
		return ""
	}
	s := prefix
	if isTail {
		s += "└──"
		prefix += "    "
	} else {
		s += "├──"
		prefix += "│   "
	}
	// Add identifier here
	// if n.isBlack() {
	// 	s += "B:"
	// } else {
	// 	s += "R:"
	// }
	// if n.parent != nil {
	// 	s += n.parent.keyString() + "->"
	// }
	s += n.keyString() + n.valString() + "\n"
	s += n.right.string(prefix, false)
	s += n.left.string(prefix, true)

	return s
}

func (n *node) keyString() string {
	if n == nil {
		return ""
	}
	return printutil.String(n.key)
}

func (n *node) valString() string {
	if n == nil {
		return ""
	}
	return fmt.Sprintf("%v", n.val)
}

func (n *node) printRoot() {
	for n.parent != nil {
		n = n.parent
	}
	fmt.Println(n)
}
