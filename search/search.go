package search

// A CompareResult is returned from a Compare query.
type CompareResult int

// Compare result constant
const (
	Less CompareResult = iota
	Equal
	Greater
	Invalid
)

// NegativeInf is a Comparable type which is less than everything (including itself!)
// (the less than itself part might change!)
type NegativeInf struct{}

// Compare on NegativeInf returns Less.
func (ni NegativeInf) Compare(c interface{}) CompareResult {
	return Less
}

// Inf is a Comparable type which is greater than everything (including itself)
type Inf struct{}

// Compare on Inf returns Greater.
func (i Inf) Compare(c interface{}) CompareResult {
	return Greater
}

// Comparable types can be compared to arbitrary
// elements and will return a CompareResult following.
// They are intended for use as search keys.
type Comparable interface {
	Compare(interface{}) CompareResult
}

// Equalable types can be compared to one another
// and will return a boolean if they are the same.
type Equalable interface {
	Equals(Equalable) bool
}

// Nil Equalables are equal to any value.
// Insert them in deletion queries in order to
// delete arbitrary elements.
type Nil struct{}

// Equals on Nil always returns true.
func (n Nil) Equals(Equalable) bool {
	return true
}

// Node types can be stored in modifiable search types.
type Node interface {
	Key() Comparable
	Val() Equalable
}

// Searchable types can be searched, with float64 keys pointing to
// arbitrary values.
type Searchable interface {
	Search(interface{}) (bool, interface{})
	SearchUp(interface{}) (Comparable, interface{})
	SearchDown(interface{}) (Comparable, interface{})
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
