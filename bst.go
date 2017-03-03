package compgeo

// BST Type enumerator
const (
	AVLTreeType = iota
	RBTreeType
	TTreeType
	SplayTreeType
	TangoTreeType
)

type Node interface {
	Key() float64
	Val() interface{}
}

type BST interface {
	Size() int
	Insert(Node) error
	Delete(Node) error
	Search(float64) (bool, interface{})
	Traverse() []interface{}
	ToPersistent() PersistentBST
}

// Hypothetical: All a PersistentBST does
// is return a pointer to a BST at a given instant
type PersistentBST interface {
	BST
	AtInstant(float64) BST
	MinInstant() float64
	MaxInstant() float64
	SetInstant(float64)
}

func NewPersistentRBTree(bstType int) (t PersistentBST) {
	switch bstType {
	case RBTreeType:
		t = NewRBTree().ToPersistent()
	case TTreeType:
		t = NewTTree().ToPersistent()
	case SplayTreeType:
		t = NewSplayTree().ToPersistent()
	case TangoTreeType:
		t = NewTangoTree().ToPersistent()
	default:
		fallthrough
	case AVLTreeType:
		t = NewAVLTree().ToPersistent()
	}
	return
}
