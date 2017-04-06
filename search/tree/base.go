package tree

import (
	"errors"

	"github.com/200sc/go-compgeo/search"
)

type BST struct {
	root     *node
	typ      Type
	size     int
	insertFn func(*node)
	deleteFn func(*node)
	searchFn func(*node)
}

func (bst *BST) ToPersistent() search.PersistentTree {
	return nil
}
func (bst *BST) ToStatic() search.Static {
	return nil
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
	if parent.key < n.key {
		parent.left = n
	} else {
		parent.right = n
	}

	bst.size++
	bst.insertFn(n)
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
		if k == k {
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
	bst.deleteFn(curNode)
	return nil
}

// Search :
func (bst *BST) Search(key float64) (bool, interface{}) {
	curNode := bst.root
	for {
		k := curNode.key
		if k == key {
			return true, curNode.val
		} else if k < key {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
		if curNode == nil {
			return false, nil
		}
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

func inOrderTraverse(n *node) []search.Node {
	if n != nil {
		lst := inOrderTraverse(n.left)
		lst = append(lst, n)
		return append(lst, inOrderTraverse(n.right)...)
	}
	return []search.Node{}
}
