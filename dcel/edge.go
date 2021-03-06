package dcel

import (
	"fmt"

	compgeo "github.com/200sc/go-compgeo"
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

// Len returns the number of discrete points
// defined on an edge, in normal cases, 2
func (e *Edge) Len() int {
	if e == nil {
		return 0
	}
	if e.Twin == nil {
		return 1
	}
	return 2
}

// At returns either e.Origin or e.Twin.Origin
// for 0 or 1
func (e *Edge) At(i int) geom.Dimensional {
	if i == 0 {
		return e.Origin
	} else if i == 1 {
		return e.Twin.Origin
	}
	panic("At exceeded dimensions on edge")
}

// Set sets the value behind a point on e
// to a given point
func (e *Edge) Set(i int, d geom.Dimensional) geom.Spanning {
	if i == 0 {
		e.Origin.Point = d.(geom.Point)
	} else if i == 1 {
		e.Twin.Origin.Point = d.(geom.Point)
	}
	return e
}

// String converts an edge into a string.
func (e *Edge) String() string {
	if e == nil {
		return "<nil>"
	}
	s := ""
	s += fmt.Sprintf("%v", e.Origin) + "->" + fmt.Sprintf("%v", e.Twin.Origin)
	return s
}

// SetTwin is shorthand for two twin assignments.
func (e *Edge) SetTwin(e2 *Edge) {
	if e != nil {
		e.Twin = e2
	}
	if e2 != nil {
		e2.Twin = e
	}
}

// SetPrev is shorthand for a prev and next assignment.
func (e *Edge) SetPrev(e2 *Edge) {
	if e != nil {
		e.Prev = e2
	}
	if e2 != nil {
		e2.Next = e
	}
}

// SetNext is shorthand for a next and prev assignment.
func (e *Edge) SetNext(e2 *Edge) {
	e2.SetPrev(e)
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
func (e *Edge) FullEdge() (geom.FullEdge, error) {
	e2 := e.Twin
	if e2 == nil {
		return geom.FullEdge{}, compgeo.BadEdgeError{}
	}
	return geom.FullEdge{
		e.Origin.Point,
		e2.Origin.Point}, nil
}

// High returns whichever point on e is higher
// in dimension d
func (e *Edge) High(d int) geom.Dimensional {
	if e == nil {
		return nil
	}
	if e.Twin == nil {
		return e.Origin
	}
	if e.Twin.Val(d) < e.Origin.Val(d) {
		return e.Origin
	}
	return e.Twin.Origin
}

// Low returns whichever point on e is lower
// in dimension d
func (e *Edge) Low(d int) geom.Dimensional {
	if e == nil {
		return nil
	}
	if e.Twin == nil {
		return e.Origin
	}
	if e.Twin.Val(d) < e.Origin.Val(d) {
		return e.Twin.Origin
	}
	return e.Origin
}

// Mid2D returns the midpoint of an Edge
func (e *Edge) Mid2D() (geom.Point, error) {
	if e == nil {
		return geom.Point{}, compgeo.BadEdgeError{}
	}
	t := e.Twin
	if t == nil {
		return geom.Point{}, compgeo.BadEdgeError{}
	}
	return e.Origin.Mid2D(t.Origin), nil
}

// Compare allows Edge to satisfy search
// interfaces for placement in BSTs.
func (e *Edge) Compare(i interface{}) search.CompareResult {
	switch c := i.(type) {
	case geom.Point:
		return geom.VerticalCompare(c, e)
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
		return false, compgeo.BadEdgeError{}
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
			return false, compgeo.BadEdgeError{}
		}
		if lowest.Origin.Greater2D(e.Origin).Eq(lowest.Origin) {
			lowest = e
		}
		e = e.Next
	}
	if lowest.Prev == nil {
		return false, compgeo.BadEdgeError{}
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
		e.Next, e.Prev = e.Prev, e.Next
		e.Twin.Next, e.Twin.Prev = e.Twin.Prev, e.Twin.Next
		e.Origin, e.Twin.Origin = e.Twin.Origin, e.Origin
		e.Origin.OutEdge = e.Origin.OutEdge.Twin
		e = e.Prev
		if e == start {
			break
		}
	}
}

// PointAt returns the point at a given position on some
// d dimension along this edge. I.E. for d = 0, v = 5,
// if this edge was represented as y = mx + b, this would
// return y = m*5 + b.
func (e *Edge) PointAt(d int, v float64) (geom.Point, error) {
	e1 := e.Origin.Point
	e2 := e.Twin.Origin.Point
	if e1[d] > e2[d] {
		e1, e2 = e2, e1
	}
	if v < e1[d] || v > e2[d] {
		return geom.Point{}, compgeo.RangeError{}
	}
	v -= e1[d]
	span := e2[d] - e1[d]
	portion := v / span
	p := geom.Point{}
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

// Val redirects to Origin.Val
func (e *Edge) Val(d int) float64 {
	return e.Origin.Val(d)
}

// D redirects to Origin.D
func (e *Edge) D() int {
	return e.Origin.D()
}

// Eq redirects to Origin.Eq
func (e *Edge) Eq(e2 geom.Spanning) bool {
	if e2.Len() != e.Len() {
		return false
	}
	for i := 0; i < e.Len(); i++ {
		if !e.At(i).Eq(e2.At(i)) {
			return false
		}
	}
	return true
}

// AllEdges on an edge is equivalent to e.Origin.AllEdges,
// which actually calls this instead of the other way
// around because that involves less code duplciation.
func (e *Edge) AllEdges() []*Edge {
	if e == nil {
		return []*Edge{}
	}
	edges := make([]*Edge, 1)
	edges[0] = e
	edge := e.Twin.Next
	for edge != e {
		edges = append(edges, edge)
		edge = edge.Twin.Next
	}
	return edges
}

// EdgeChain returns the edges of e as a chain
// until e.Next = e again.
func (e *Edge) EdgeChain() []*Edge {
	out := []*Edge{e}
	for e2 := e.Next; e2 != e; e2 = e2.Next {
		out = append(out, e2)
	}
	return out
}

// FindSharedPoint returns some value, if one exists,
// in dimension d, that has a defined point on both
// e and e2.
// This function is deterministic. It does not return
// a -random- value, if multiple values are valid.
func (e *Edge) FindSharedPoint(e2 *Edge, d int) (float64, error) {
	eLow := e.Low(d)
	e2Low := e2.Low(d)
	eHigh := e.High(d)
	e2High := e2.High(d)
	// The two edges span the same distance case
	if geom.F64eq(eLow.Val(d), e2Low.Val(d)) &&
		geom.F64eq(eHigh.Val(d), e2High.Val(d)) {
		return ((eHigh.Val(d) - eLow.Val(d)) / 2) + eLow.Val(d), nil
	}

	// No point exists case:
	// if e.High is lower than e2.Low
	// or if e.Low is higher than e2.High
	if eHigh.Val(d) < e2Low.Val(d) ||
		eLow.Val(d) > e2High.Val(d) {
		return 0, compgeo.BadEdgeError{}
	}
	if eHigh.Val(d) > e2High.Val(d) &&
		eLow.Val(d) < e2Low.Val(d) {
		// E completely contains e2 in this dimension
		mid, err := e2.Mid2D()
		if err != nil {
			return 0, err
		}
		return mid.Val(d), nil
	}
	if eHigh.Val(d) < e2High.Val(d) &&
		eLow.Val(d) > e2Low.Val(d) {
		// E2 completely contains e in this dimension
		mid, err := e.Mid2D()
		if err != nil {
			return 0, err
		}
		return mid.Val(d), nil
	}
	r1 := eHigh.Val(d) - e2Low.Val(d)
	r2 := e2High.Val(d) - eLow.Val(d)
	// fmt.Println("Found shared point!")
	// fmt.Println(e, e2)
	if r1 < r2 {
		//fmt.Println(r1/2 + e2Low.Val(d))
		return r1/2 + e2Low.Val(d), nil
	}
	//fmt.Println(r2/2 + eLow.Val(d))
	return r2/2 + eLow.Val(d), nil
}

// PointAlong returns the point some percent along
// e, increasing on dimension d. If percent <= 0, returns e.Low(d).
// if percent >= 1, returns e.High(d).
func (e *Edge) PointAlong(d int, pcnt float64) geom.D3 {
	p1 := e.Low(d).(geom.D3)
	p2 := e.High(d).(geom.D3)
	if pcnt <= 0 {
		return p1
	}
	if pcnt >= 1 {
		return p2
	}
	return geom.NewPoint(
		p1.X()+(p2.X()-p1.X())*pcnt,
		p1.Y()+(p2.Y()-p1.Y())*pcnt,
		p1.Z()+(p2.Z()-p1.Z())*pcnt)
}

// Copy returns a copy of edge e
func (e *Edge) Copy() *Edge {
	e2 := new(Edge)
	e2.Origin = e.Origin
	e2.Face = e.Face
	e2.Prev = e.Prev
	e2.Next = e.Next
	e2.Twin = e.Twin
	return e2
}
