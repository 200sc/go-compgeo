package tree

import (
	"errors"
	"fmt"
	"math"

	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree/static"
)

type fnSet struct {
	insertFn func(*node) *node
	deleteFn func(*node) *node
	searchFn func(*node) *node
}

func nopNode(n *node) *node {
	return nil
}

// BST is a generic binary search tree implementation.
// BST relies on the idea that numberous BST types are
// implicitly the same, but with unique functions to update
// their balance after each insert, delete, or search (sometimes)
// operation.
type BST struct {
	*fnSet
	root *node
	// Because the size of a bst is something someone might want
	// to query quickly, we raise it to the top instead of making
	// it a tree-wide count-up.
	size int
}

func (bst *BST) isValid() bool {
	ok, _, _ := bst.root.isValid()
	return ok
}

// ToPersistent converts this BST into a Persistent BST.
func (bst *BST) ToPersistent() search.DynamicPersistent {
	pbst := new(PersistentBST)
	pbst.instant = math.MaxFloat64 * -1
	pbst.instants = []BSTInstant{{BST: bst, instant: pbst.instant}}
	return pbst
}

// ToStatic on a BST figures out where all nodes
// would exist in an array structure, then constructs
// an array with a length of the maximum index found.
//
// If static stays in its own package this presents
// a potential import cycle-- or else all of static's
// tests need to exist outside of static, as it can't
// create an instance of a staticBST by itself.
func (bst *BST) ToStatic() search.Static {
	m, maxIndex := bst.root.staticTree(make(map[int]*static.Node), 1)
	staticBst := make(static.BST, maxIndex+1)
	for k, v := range m {
		staticBst[k] = v
	}
	return &staticBst
}

// Size :
func (bst *BST) Size() int {
	return bst.size
}

// Insert :
func (bst *BST) Insert(inNode search.Node) error {
	n := new(node)
	n.key = inNode.Key()
	n.val = inNode.Val()
	n.payload = red
	var parent *node
	curNode := bst.root
	for {
		if curNode == nil {
			break
		}
		parent = curNode
		if curNode.key > n.key {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
		// We do not do any sort of checking for duplicates in this type.
		// This means the same key and value pair, or two values with the
		// same key can both be in this tree.
		// Todo: if we need the type, create treeSet types which
		// do nothing on duplicates being added.
	}
	// curNode == nil
	n.parent = parent
	if parent != nil {
		if parent.key > n.key {
			parent.left = n
		} else {
			parent.right = n
		}
		// if parent == nil and curNode == nil,
		// this bst is empty.
	} else {
		n.payload = black
		bst.root = n
	}

	bst.size++
	bst.updateRoot(bst.insertFn(n))
	return nil
}

// Delete :
// Because we allow duplicate keys,
// because real data has duplicate keys,
// we require you specify what you want to delete
// at the given key or nil if you know for sure that
// there is only one value with the given key (or
// do not care what is deleted).
func (bst *BST) Delete(n search.Node) error {
	curNode := bst.root
	v := n.Val()
	k := n.Key()
	for {
		k2 := curNode.key
		if k2 == k {
			if v == nil {
				break
			}
			for curNode.val != v {
				// We're only going to find keys that are the same
				// as this key in this key's right descendants
				curNode = curNode.right
				if curNode == nil {
					return errors.New("Value not found")
				}
			}
			break
		} else if k2 > k {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
		if curNode == nil {
			return errors.New("Key not found")
		}
	}
	bst.size--
	bst.updateRoot(bst.deleteFn(curNode))
	return nil
}

// Search :
func (bst *BST) Search(key float64) (bool, interface{}) {
	curNode := bst.root
	var k float64
	for curNode != nil {
		k = curNode.key
		if k == key {
			break
		} else if k > key {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
	}
	if curNode != nil {
		bst.updateRoot(bst.searchFn(curNode))
		return true, curNode.val
	}
	return false, nil
}

func (bst *BST) search(key float64) (*node, bool) {
	curNode := bst.root
	var k float64
	var parent *node
	for curNode != nil {
		k = curNode.key
		parent = curNode
		if k == key {
			break
		} else if k > key {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
	}
	if curNode != nil {
		return curNode, true
	}
	return parent, false
}

// SearchUp performs a search, and rounds up to the nearest
// existing key if no node of the query key exists.
func (bst *BST) SearchUp(key float64) interface{} {
	n, ok := bst.search(key)
	// The tree is empty
	if n == nil {
		return nil
	}
	if ok {
		return n.val
	}
	v := n.successor()
	if v == nil || ((v.key > n.key) && (n.key > key)) {
		return n.val
	}
	return v.val
}

// SearchDown acts as SearchUp, but rounds down.
func (bst *BST) SearchDown(key float64) interface{} {
	n, ok := bst.search(key)
	if ok {
		return n.val
	}
	v := n.predecessor()
	if v == nil || ((v.key < n.key) && (n.key < key)) {
		return n.val
	}
	return v.val
}

func (bst *BST) updateRoot(n *node) {
	if n != nil {
		bst.root = n
		return
	}
	if bst.size == 0 {
		bst.root = nil
	}
	if bst.root == nil {
		return
	}
	for bst.root.parent != nil {
		bst.root = bst.root.parent
	}
}

// InOrderTraverse :
// There are multiple ways to traverse a tree.
// The most useful of these is the in-order traverse,
// and that's what we provide here.
// Other traversal methods can be added as needed.
func (bst *BST) InOrderTraverse() []search.Node {
	return inOrderTraverse(bst.root)
}

func (bst *BST) copy() *BST {
	newBst := new(BST)
	newBst.root = bst.root.copy()
	newBst.fnSet = bst.fnSet
	return newBst
}

func (bst *BST) String() string {
	s := bst.root.string("", true)
	if s == "" {
		return "<Empty BST>\n"
	}
	return s
}

func findCycle(bst *BST) error {
	seen := make(map[float64]map[float64]bool)
	return bst.root.findCycle(seen)
}

// findCycle will mis-report duplicate input nodes as cycles.
func (n *node) findCycle(seen map[float64]map[float64]bool) error {
	if n == nil {
		return nil
	}
	if v, ok := seen[n.key]; ok {
		if _, ok = v[n.val.(float64)]; ok {
			fmt.Println(n)
			return errors.New("Cycle found")
		}
		seen[n.key][n.val.(float64)] = true
	} else {
		seen[n.key] = make(map[float64]bool)
		seen[n.key][n.val.(float64)] = true
	}
	err := n.left.findCycle(seen)
	if err != nil {
		return err
	}
	return n.right.findCycle(seen)
}
