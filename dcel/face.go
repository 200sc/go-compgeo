package dcel

import (
	"sort"

	"github.com/200sc/go-compgeo/geom"
)

// A Face points to the edges on its inner and
// outer portions. Any given face may have either
// of these values be nil, but never both.
type Face struct {
	Inner, Outer *Edge
}

// NewFace returns a null-initialized Face.
func NewFace() *Face {
	return &Face{}
}

// Vertices wraps around a face and
// finds all vertices that border it.
func (f *Face) Vertices() []*Vertex {
	// Inner is not populated by anything as of this writing.

	pts := []*Vertex{}
	e := f.Outer
	for e != nil && e.Next != f.Outer {
		pts = append(pts, e.Origin)
		e = e.Next
	}
	if e != nil {
		pts = append(pts, e.Origin)
	}
	return pts
}

// Contains returns whether a point lies inside f.
// We cannot assume that f is convex, or anything
// besides some polygon. That leaves us with a rather
// complex form of PIP--
func (f *Face) Contains(p geom.D2) bool {
	x := p.X()
	y := p.Y()
	contains := false
	bounds := f.Bounds()
	min := bounds.At(0).(geom.D2)
	max := bounds.At(1).(geom.D2)
	if x < min.Val(0) || x > max.Val(0) ||
		y < min.Val(1) || y > max.Val(1) {
		return contains
	}

	e1 := f.Outer.Prev
	e2 := f.Outer
	for {
		if (e2.Y() > y) != (e1.Y() > y) {
			if x < (e1.X()-e2.X())*(y-e2.Y())/(e1.Y()-e2.Y())+e2.X() {
				contains = !contains
			}
		}
		e1 = e1.Next
		e2 = e2.Next
		if e1 == f.Outer.Prev {
			break
		}
	}
	return contains
}

// VerticesSorted returns this face's vertices sorted in dimensions ds.
// Example: to get points sorted by x, use with (0)
//          to get points sorted by y, breaking ties
//             on lesser x, use with (1,0).
// --This has different behavior than DCEL.VerticesSorted!
// it does not return indices but direct vertex pointers.
func (f *Face) VerticesSorted(ds ...int) []*Vertex {
	pts := f.Vertices()
	sort.Slice(pts, func(i, j int) bool {
		for _, d := range ds {
			v1 := pts[i].Val(d)
			v2 := pts[j].Val(d)
			if v1 != v2 {
				return v1 < v2
			}
		}
		return false
	})
	return pts
}

// Encloses returns whether f completey enwraps f2
// Doing this check legitimately would
// be costly and complex. We assume, right now,
// that we already -know- that either f encloses f2
// or f2 encloses f.
// If this is true, if one of them has a point higher
// than the other, that one is the encloser.
func (f *Face) Encloses(f2 *Face) bool {
	v1 := f.Vertices()
	if len(v1) == 0 {
		// f is the outer face
		return true
	}
	v2 := f2.Vertices()
	if len(v2) == 0 {
		return false
	}
	fMax := geom.NegInf
	f2Max := geom.NegInf
	for _, v := range v1 {
		if v.Y() > fMax {
			fMax = v.Y()
		}
	}
	for _, v := range v2 {
		if v.Y() > f2Max {
			f2Max = v.Y()
		}
	}
	return fMax > f2Max
}

const (
	// OUTER_FACE is used to represent the infinite space
	// around the outer edge(s) of a DCEL.
	OUTER_FACE = 0
)
