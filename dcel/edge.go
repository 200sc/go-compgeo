package dcel

import (
	"fmt"

	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search"
)

// An Edge represents an edge within a DCEL,
// specifically a half edge, which maintains
// references to it's origin vertex, the face
// it bounds, the half edge sharing its space
// bounding its adjacent face, and the previous
// and following edges which bound its face.
type Edge struct {
	// Origin is the vertex this edge starts at
	Origin *Vertex
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

// NewEdge returns a null-initialized Edge.
func NewEdge() *Edge {
	return &Edge{}
}

// String converts an edge into a string.
func (e *Edge) String() string {
	s := ""
	s += fmt.Sprintf("%v", e.Origin) + "->" + fmt.Sprintf("%v", e.Twin.Origin)
	return s
}

// EdgeTwin can obtain a given edge index's twin
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

// FullEdge returns this edge with its twin in the form of its
// two vertices
func (e *Edge) FullEdge() (FullEdge, error) {
	e2 := e.Twin
	if e2 == nil {
		return FullEdge{}, BadEdgeError{}
	}
	return FullEdge{
		e.Origin.Point,
		e2.Origin.Point}, nil
}

// Mid2D returns the midpoint of an Edge
func (e *Edge) Mid2D() (*Point, error) {
	if e == nil {
		return nil, BadEdgeError{}
	}
	t := e.Twin
	if t == nil {
		return nil, BadEdgeError{}
	}
	return e.Origin.Mid2D(t.Origin), nil
}

// Compare allows Edge to satisfy search
// interfaces for placement in BSTs.
func (e *Edge) Compare(i interface{}) search.CompareResult {
	switch c := i.(type) {
	case Point:
		return c.VerticalCompare(e)
	case *Point:
		return c.VerticalCompare(e)
	default:
		return search.Invalid
	}
}

// IsClockwise returns whether a given set of
// edges is clockwise or not.
// Method credit: lhf on stackOverflow
// https://math.stackexchange.com/questions/340830
func (e *Edge) IsClockwise() (bool, error) {
	if e == nil {
		return false, BadEdgeError{}
	}
	start := e
	lowest := e
	// Find the highest, rightmost point.
	// We find the highest because the axes
	// in this system are flipped so y increases
	// going downward. Ultimately as long as we
	// are consistent with one approach this does
	// not change anything.
	e = e.Next
	for e != start {
		if e == nil || e.Next == nil {
			return false, BadEdgeError{}
		}
		if lowest.Origin.Greater2D(e.Origin).Eq(lowest.Origin) {
			lowest = e
		}
		e = e.Next
	}
	if lowest.Prev == nil {
		return false, BadEdgeError{}
	}
	p := lowest.Prev.Origin.Point
	c := lowest.Origin.Point
	n := lowest.Next.Origin.Point

	cross := (p[0] * c[1]) - (c[0] * p[1]) +
		(p[1] * n[0]) - (p[0] * n[1]) +
		(c[0] * n[1]) - (n[0] * c[1])

	// We assume the points are not colinear,
	// as they must not be. If they were,
	// one of lowest's neighbors would be
	// higher than lowest.
	if cross > 0 {
		return true, nil
	}
	return false, nil
}

// Flip converts edge and all that share a
// face with edge from counterclockwise to clockwise
// or vice versa
func (e *Edge) Flip() {
	start := e
	for {
		//fmt.Println("Flipping!")
		e.Next, e.Prev = e.Prev, e.Next
		e.Twin.Next, e.Twin.Prev = e.Twin.Prev, e.Twin.Next
		e.Origin, e.Twin.Origin = e.Twin.Origin, e.Origin
		e.Origin.OutEdge = e.Origin.OutEdge.Twin
		e = e.Prev
		if e == start {
			break
		}
	}
	// for {
	// 	fmt.Println("E:", e)
	// 	fmt.Println("Twin", e.Twin)
	// 	e = e.Next
	// 	if e == start {
	// 		break
	// 	}
	// }
}

// PointAt returns the point at a given position on some
// d dimension along this edge. I.E. for d = 0, v = 5,
// if this edge was represented as y = mx + b, this would
// return y = m*5 + b.
func (e *Edge) PointAt(d int, v float64) (*Point, error) {
	e1 := e.Origin.Point
	e2 := e.Twin.Origin.Point
	if e1[d] > e2[d] {
		e1, e2 = e2, e1
	}
	if v < e1[d] || v > e2[d] {
		fmt.Println(v, e1[d], e2[d])
		return nil, RangeError{}
	}
	v -= e1[d]
	span := e2[d] - e1[d]
	portion := v / span
	p := new(Point)
	for i := 0; i < len(p); i++ {
		if i == d {
			p[i] = v
		} else {
			p[i] = e1[i] + (portion * (e2[i] - e1[i]))
		}
	}
	return p, nil
}

// Y redirects to Origin.Y
func (e *Edge) Y() float64 {
	return e.Origin.Y()
}

// X redirects to Origin.X
func (e *Edge) X() float64 {
	return e.Origin.X()
}

// Z redirects to Origin.Z
func (e *Edge) Z() float64 {
	return e.Origin.Z()
}

func (e *Edge) Val(d int) float64 {
	return e.Origin.Val(d)
}
func (e *Edge) D() int {
	return e.Origin.D()
}
func (e *Edge) Eq(e2 geom.Dimensional) bool {
	return e.Origin.Eq(e2)
}

// AllEdges on an edge is equivalent to e.Origin.AllEdges,
// which actually calls this instead of the other way
// around because that involves less code duplciation.
func (e *Edge) AllEdges() []*Edge {
	edges := make([]*Edge, 1)
	edges[0] = e
	edge := e.Twin.Next
	for edge != e {
		edges = append(edges, edge)
		edge = edge.Twin.Next
	}
	return edges
}
