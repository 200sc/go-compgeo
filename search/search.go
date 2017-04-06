package search

// Node types can be stored in modifiable search types.
type Node interface {
	Key() float64
	Val() interface{}
}

// Searchable types can be searched, with float64 keys pointing to
// arbitrary values.
type Searchable interface {
	Search(float64) (bool, interface{})
}

// Traversable types can produce lists of elements
// in some given order from within them.
type Traversable interface {
	InOrderTraverse() []Node
}

// Sizable types can count the number of elements within them.
type Sizable interface {
	Size() int
}

// Static types can be searched and traversed,
// but not modified. The benefit of using a static
// type is that they should be faster to query than
// a dynamic type.
type Static interface {
	Sizable
	Searchable
	Traversable
}

// Modifiable types implicitly are static types with
// additional functionality to change their contents,
// and being static types, should be convertible back to static types.
type Modifiable interface {
	Static
	Insert(Node) error
	Delete(Node) error
	ToStatic() Static
}
