package tree

import (
	"fmt"
	"strconv"

	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree/static"
)

type node struct {
	// eventually key should be a comparable interface
	// but that would probably poorly effect performance
	key float64
	val interface{}
	// Each tree type might have a different payload on each node
	// a good example of this is RED or BLACK on RBtrees.
	payload interface{}

	left, right, parent *node
}

func (n *node) Key() float64 {
	return n.key
}

func (n *node) Val() interface{} {
	return n.val
}

func (n *node) copy() *node {
	if n == nil {
		return nil
	}
	cp := new(node)
	left := n.left.copy()
	right := n.right.copy()

	cp.left = left
	if cp.left != nil {
		cp.left.parent = cp
	}
	cp.right = right
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

func parent_sibling(n, p *node) *node {
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

func (n *node) leftRotate() {
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
	}
	r.left = n
	n.parent = r
}

func (n *node) rightRotate() {
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
	}
	l.right = n
	n.parent = l
}

func (n *node) swap(n2 *node) {
	// Reset the second node's parent reference
	if n2.parent != nil {
		if n2.parent.left == n2 {
			n2.parent.left = n
		} else {
			n2.parent.right = n
		}
	}
	n2parent := n2.parent
	n2left := n2.left
	n2right := n2.right

	// Point deleted node's children to the lifted node,
	// and vice versa.
	n2.left = n.left
	n2.right = n.right
	if n2.left != nil {
		n2.left.parent = n2
	}
	if n2.right != nil {
		n2.right.parent = n2
	}

	n2.parent = n.parent
	// Repeat for the first node
	if n2.parent != nil {
		if n2.parent.left == n {
			n2.parent.left = n2
		} else {
			n2.parent.right = n2
		}
	}

	n.left = n2left
	n.right = n2right
	if n.left != nil {
		n.left.parent = n
	}
	if n.right != nil {
		n.right.parent = n
	}

	n.parent = n2parent
}

func (n *node) staticTree(m map[int]static.Node, i int) (map[int]static.Node, int) {
	if n == nil {
		return m, 0
	}
	m[i] = static.NewNode(n.key, n.val)
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
	if n.isBlack() {
		s += "B:"
	} else {
		s += "R:"
	}
	if n.parent != nil {
		s += keyString(n.parent.key) + "->"
	}
	s += keyString(n.key) + "\n"
	s += n.left.string(prefix, false)
	s += n.right.string(prefix, true)

	return s
}

func keyString(k float64) string {
	return strconv.FormatFloat(k, 'f', -1, 64)
}

func (n *node) keyString() string {
	if n == nil {
		return ""
	}
	return keyString(n.key)
}

func (n *node) printRoot() {
	for n.parent != nil {
		n = n.parent
	}
	fmt.Println(n)
}
