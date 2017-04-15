package dcel

import (
	"fmt"
	"image/color"
)

// A Face points to the edges on its inner and
// outer portions. Any given face may have either
// of these values be nil, but never both.
//
// We are using Outer and Inner completely wrong at time of writing.
// Outer is unused when Inner should be.
type Face struct {
	Outer, Inner *Edge
	Color        color.Color
}

// NewFace returns a null-initialized Face.
func NewFace() *Face {
	return &Face{Color: color.RGBA{0, 255, 255, 255}}
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

func (f *Face) CorrectDirectionality() {
	// Inners need to be going CC
	// Outers need to be going Clockwise

	clock, err := f.Inner.IsClockwise()
	if err == nil && clock {
		f.Inner.Flip()
	} else {
		fmt.Println(err, clock)
	}
	clock, err = f.Outer.IsClockwise()
	if err == nil && !clock {
		f.Outer.Flip()
	}
}

// Encloses returns whether f compleletly enwraps
// f2.
func (f *Face) Encloses(f2 *Face) bool {
	// Doing this check legitimately would
	// be costly and complex. We assume, right now,
	// that we already -know- that either f encloses f2
	// or f2 encloses f.
	// If this is true, if one of them has a point higher
	// than the other, that one is the encloser.
	return (f.Max(1) > f2.Max(1))
}

func (f *Face) Max(i int) (x float64) {
	e := f.Inner
	max := e.Origin
	for e.Next != f.Inner {
		e = e.Next
		if e.Origin[i] > max[i] {
			max = e.Origin
		}
	}
	return max[i]
}

const (
	// OUTER_FACE is used to represent the infinite space
	// around the outer edge of a DCEL.
	OUTER_FACE = 0
)
