package compgeo

type Node interface {
	Key() float64
	Val() interface{}
}

type Searchable interface {
	Search(float64) (bool, interface{})
}

type Traverseable interface {
	InOrderTraverse() []Node
}

type StaticSearchTree interface {
	Searchable
	Traverseable
	IsNil(int) bool
}

type SearchTree interface {
	StaticSearchTree
	Insert(Node) error
	Delete(Node) error
	ToPersistent() PersistentBST
	ToStatic() StaticSearchTree
	Size() int
}

type PersistentBST interface {
	SearchTree
	AtInstant(float64) SearchTree
	MinInstant() float64
	MaxInstant() float64
	SetInstant(float64)
}

type TreeType int

const (
	AVLTreeType TreeType = iota
	RBTreeType
	TTreeType
	SplayTreeType
	TangoTreeType
)
