package dcel

import "fmt"

// A Face points to the edges on its inner and
// outer portions. Any given face may have either
// of these values be nil, but never both.
//
// We are using Outer and Inner completely wrong at time of writing.
// Outer is unused when Inner should be.
type Face struct {
	Outer, Inner *Edge
}

// NewFace returns a null-initialized Face.
func NewFace() *Face {
	return &Face{}
}

// Vertices wraps around a face and
// finds all vertices that border it.
func (f *Face) Vertices() []Point {
	// Outer is not populated by anything as of this writing.

	pts := []Point{}
	e := f.Inner
	for e != nil && e.Next != f.Inner {
		pts = append(pts, *e.Origin)
		e = e.Next
	}
	if e != nil {
		pts = append(pts, *e.Origin)
	}
	return pts
}

// Contains2D is an interface which satisfies
// the functions needed to determine if
// a one-dimensional element exists on a face
type Contains2D interface {
	X() float64
	Y() float64
}

// Contains returns whether a point lies inside f.
// We cannot assume that f is convex, or anything
// besides some polygon. That leaves us with a rather
// complex form of PIP--
func (f *Face) Contains(p Contains2D) bool {
	x := p.X()
	y := p.Y()
	contains := false
	bounds := f.Bounds()
	fmt.Println("Face bounds", bounds)
	if x < bounds.Min.X() || x > bounds.Max.X() ||
		y < bounds.Min.Y() || y > bounds.Max.Y() {
		return contains
	}
	fmt.Println("Point lied in bounds")

	e1 := f.Inner.Prev
	e2 := f.Inner
	for {
		if (e2.Y() > y) != (e1.Y() > y) { // Three comparisons
			if x < (e1.X()-e2.X())*(y-e2.Y())/(e1.Y()-e2.Y())+e2.X() { // One Comparison, Four add/sub, Two mult/div
				contains = !contains
			}
		}
		e1 = e1.Next
		e2 = e2.Next
		if e1 == f.Inner.Prev {
			break
		}
	}
	return contains
}

const (
	// OUTER_FACE is used to represent the infinite space
	// around the outer edge(s) of a DCEL.
	OUTER_FACE = 0
)
