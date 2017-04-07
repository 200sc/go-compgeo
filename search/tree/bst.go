package tree

import (
	"errors"
	"fmt"

	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree/static"
)

type fnSet struct {
	insertFn func(*node)
	deleteFn func(*node)
	searchFn func(*node)
}

type BST struct {
	*fnSet
	root *node
	size int
}

func (bst *BST) ToPersistent() search.DynamicPersistent {
	return nil
}

// ToStatic on a BST figures out where all nodes
// would exist in an array structure, then constructs
// an array with a length of the maximum index found.
func (bst *BST) ToStatic() search.Static {
	m, maxIndex := bst.root.staticTree(make(map[int]static.Node), 1)
	staticBst := make(static.BST, maxIndex)
	for k, v := range m {
		staticBst[k] = v
	}
	return staticBst
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
		if curNode.key < n.key {
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
		if parent.key < n.key {
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
	bst.insertFn(n)
	bst.updateRoot()
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
		} else if k2 < k {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
		if curNode == nil {
			return errors.New("Key not found")
		}
	}
	bst.size--
	bst.deleteFn(curNode)
	bst.updateRoot()
	return nil
}

// Search :
func (bst *BST) Search(key float64) (bool, interface{}) {
	curNode := bst.root
	for {
		if curNode == nil {
			break
		}
		k := curNode.key
		if k == key {
			break
		} else if k < key {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
	}
	if curNode != nil {
		bst.searchFn(curNode)
		bst.updateRoot()
		return true, curNode.val
	}
	return false, nil
}

func (bst *BST) updateRoot() {
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
		return "<Empty BST>"
	}
	return s
}

func findCycle(bst *BST) {
	seen := make(map[float64]map[float64]bool)
	bst.root.findCycle(seen)
}

func (n *node) findCycle(seen map[float64]map[float64]bool) {
	if n == nil {
		return
	}
	if v, ok := seen[n.key]; ok {
		if _, ok = v[n.val.(float64)]; ok {
			fmt.Println(n)
			panic("Cycle found")
		} else {
			seen[n.key][n.val.(float64)] = true
		}
	} else {
		seen[n.key] = make(map[float64]bool)
		seen[n.key][n.val.(float64)] = true
	}
	n.left.findCycle(seen)
	n.right.findCycle(seen)
}
