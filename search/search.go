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
	SearchUp(float64) interface{}
	SearchDown(float64) interface{}
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
// type is that it should be faster to query than
// a dynamic type.
type Static interface {
	Sizable
	Searchable
	Traversable
}

// Dynamic types implicitly are static types with
// additional functionality to change their contents,
// and being static types, should be convertible back to static.
type Dynamic interface {
	Static
	Insert(Node) error
	Delete(Node) error
	ToStatic() Static
}

// Why isn't there a Static -> Dynamic function?
// Because Static types don't store the information required to
// make modifications to their structure, and don't have an idea of what
// information they would add to do so-- that's the job of individual
// Dynamic structures. So yes, functions like ToRBTree(Static)
// might exist, but they won't be functions on the interface itself.
