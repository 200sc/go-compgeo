package dcel

// A Face points to the edges on its inner and
// outer portions. Any given face may have either
// of these values be nil, but never both.
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

const (
	// OOUTER_FACE is used to represent the infinite space
	// around the outer edge of a DCEL.
	OUTER_FACE = 0
)
