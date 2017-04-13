package dcel

import "fmt"

// A DCELEdge represents an edge within a DCEL,
// specifically a half edge, which maintains
// references to it's origin vertex, the face
// it bounds, the half edge sharing its space
// bounding its adjacent face, and the previous
// and following edges which bound its face.
type Edge struct {
	// Origin is the vertex this edge starts at
	Origin *Point
	// Face is the index within Faces that this
	// edge wraps around
	Face *Face
	// Next and Prev are the edges following and
	// preceding this edge that also wrap around
	// Face
	Next *Edge
	Prev *Edge
	// Twin is the half edge who points to this
	// half-edge's origin, and respectively whose
	// origin this half-edge points to.
	Twin *Edge
}

func NewEdge() *Edge {
	return &Edge{}
}

func (e *Edge) String() string {
	s := ""
	s += fmt.Sprintf("%v", e.Origin) + "->" + fmt.Sprintf("%v", e.Twin.Origin)
	return s
}

// DCELEdgeTwin can obtain a given edge index's twin
// without accessing the edge itself, for index
// manipulation, or for initially setting the Twins
// in construction.
//
// Hopeful Mandate: twin edges come in pairs
// if i is even, then, i+1 is its pair,
// and otherwise i-i is its pair.
func EdgeTwin(i int) int {
	if i%2 == 0 {
		return i + 1
	}
	return i - 1
}

// FullEdge returns the ith edge in the form of its
// two vertices
func (d *DCEL) FullEdge(i int) ([2]*Point, error) {
	if i >= len(d.HalfEdges) {
		return [2]*Point{}, BadEdgeError{}
	}
	e := d.HalfEdges[i]
	if e.Twin == nil {
		return [2]*Point{}, BadEdgeError{}
	}
	e2 := e.Twin
	return [2]*Point{
		e.Origin,
		e2.Origin}, nil
}

type BadEdgeError struct{}

func (bee BadEdgeError) Error() string {
	return "The input edge was invalid"
}
