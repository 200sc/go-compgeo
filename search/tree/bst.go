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
	//fmt.Println("Inserting ", inNode)
	n := new(node)
	n.key = inNode.Key()
	n.val = []interface{}{inNode.Val()}
	// We can't do this once we have more than RB trees wow
	n.payload = red
	var parent *node
	curNode := bst.root
	for {
		if curNode == nil {
			break
		}
		parent = curNode
		r := curNode.key.Compare(n.key)
		if r == search.Greater {
			curNode = curNode.left
		} else if r == search.Less {
			curNode = curNode.right
		} else if r == search.Equal {
			// All values of the same key are stored at the same node
			curNode.val = append(curNode.val, inNode.Val())
			bst.size++
			return nil
		} else {
			panic("Invalid types for BST operatrions")
		}
		// Todo: if we need the type, create treeSet types which
		// do nothing on duplicates being added.
	}
	// curNode == nil
	n.parent = parent
	if parent != nil {
		if parent.key.Compare(n.key) == search.Greater {
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
	//fmt.Println("Inserting for real!", n, parent)
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
	curNode, isReal := bst.search(k)
	if !isReal {
		return errors.New("Key not found")
	}
	if len(curNode.val) != 1 {
		// Otherwise if no value to delete was specified, delete
		// the first value in this list.
		if v == nil {
			curNode.val = curNode.val[len(curNode.val)-1:]
			bst.size--
			return nil
		}
		// Otherwise scan to find the value to delete.
		// If this becomes a performance hit, the user
		// should consider whether some part of the value
		// should not be encoded into the key.
		for vi := 0; vi < len(curNode.val); vi++ {
			if curNode.val[vi] == v {
				curNode.val = append(curNode.val[:vi], curNode.val[vi+1:]...)
				bst.size--
				return nil
			}
		}
		return errors.New("Value not found")
	}
	bst.size--
	bst.updateRoot(bst.deleteFn(curNode))
	return nil
}

// Search :
func (bst *BST) Search(key interface{}) (bool, interface{}) {
	curNode, isReal := bst.search(key)
	if !isReal {
		return false, nil
	}
	bst.updateRoot(bst.searchFn(curNode))
	return true, curNode.val[0]
}

func (bst *BST) search(key interface{}) (*node, bool) {
	curNode := bst.root
	var k search.Comparable
	var parent *node
	for curNode != nil {
		k = curNode.key
		parent = curNode
		r := k.Compare(key)
		if r == search.Equal {
			break
		} else if r == search.Greater {
			curNode = curNode.left
		} else if r == search.Less {
			curNode = curNode.right
		} else {
			panic("Invalid types for BST operations")
		}
	}
	if curNode != nil {
		return curNode, true
	}
	return parent, false
}

// SearchUp performs a search, and rounds up to the nearest
// existing key if no node of the query key exists.
func (bst *BST) SearchUp(key interface{}) (search.Comparable, interface{}) {
	n, ok := bst.search(key)
	// The tree is empty
	if n == nil {
		return nil, nil
	}
	if ok {
		return n.key, n.val[0]
	}
	v := n.successor()
	if v == nil ||
		((v.key.Compare(n.key) == search.Greater) &&
			(n.key.Compare(key) == search.Greater)) {
		return n.key, n.val[0]
	}
	return v.key, v.val[0]
}

// SearchDown acts as SearchUp, but rounds down.
func (bst *BST) SearchDown(key interface{}) (search.Comparable, interface{}) {
	n, ok := bst.search(key)
	if n == nil {
		return nil, nil
	}
	if ok {
		return n.key, n.val[0]
	}
	v := n.predecessor()
	if v == nil ||
		((v.key.Compare(n.key) == search.Less) &&
			n.key.Compare(key) == search.Less) {
		return n.key, n.val[0]
	}
	return v.key, v.val[0]
}

func (bst *BST) updateRoot(n *node) {
	if bst.size == 0 {
		bst.root = nil
		return
	}
	if n != nil {
		bst.root = n
		//fmt.Println("Setting root to", n)
		return
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
	newBst.size = bst.size
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
	seen := make(map[search.Comparable]bool)
	return bst.root.findCycle(seen)
}

// findCycle will mis-report duplicate input nodes as cycles.
func (n *node) findCycle(seen map[search.Comparable]bool) error {
	if n == nil {
		return nil
	}
	if _, ok := seen[n.key]; ok {
		fmt.Println(n)
		return errors.New("Cycle found")
	}
	seen[n.key] = true

	err := n.left.findCycle(seen)
	if err != nil {
		return err
	}
	return n.right.findCycle(seen)
}
