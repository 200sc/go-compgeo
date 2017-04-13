package tree

import "github.com/200sc/go-compgeo/search"

// Type represents the underlying algorithm for updating points on
// a dynamic binary search tree.
// This implementation relies on the idea that, in principle, all
// binary search trees share a lot in common (finding where to
// insert, delete, search), and any remaining details just depend
// on what specific BST type is being used.
type Type int

// TreeType enum
const (
	AVL      Type = iota
	RedBlack      // RB would probably be okay.
	Splay
	// Consider:
	// Treap?
	// Scapegoat tree?
	// TTree? <- more work than the other two
	// AA?
)

// New returns a tree as defined by the input type.
// Hypothetically, this is the only exported function in this package
func New(typ Type) search.Persistable {
	bst := new(BST)
	switch typ {
	case AVL:
		bst.fnSet = avlFnSet
	case Splay:
		bst.fnSet = splayFnSet
	default:
		fallthrough
	case RedBlack:
		bst.fnSet = rbFnSet
	}
	return bst
}