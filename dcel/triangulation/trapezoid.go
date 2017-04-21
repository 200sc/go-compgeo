package triangulation

import "github.com/200sc/go-compgeo/dcel"

// These constants refer to indices
// within trapezoids' Edges
const (
	top = iota
	bot
	left
	right
)

// These constants refer to indices
// within trapezoids' Neighbors
const (
	upright = iota
	botright
	upleft
	botleft
)

// A Trapezoid is used when contstructing a Trapezoid map,
// and contains references to its neighbor trapezoids and
// the edges that border it.
type Trapezoid struct {
	// See above indices
	Edges     [4]dcel.FullEdge
	Neighbors [4]*Trapezoid
	node      *TrapezoidNode
	face      *dcel.Face
}

// DCELEdges evaluates and returns the edges of
// a trapezoid as DCElEdges with initialized origins,
// prevs, and nexts.
//
// DCELEdges makes one assumption about our data:
// it is very unlikely that we will have an innaccurate top
// or bottom, but potentially likely we will have an innaccurate
// left or right. (We also may do away with left and right,
// as their information is stored in top and bottom)
// By this, this function adds each successive vertex that is not
// the same as the previous added vertex in order--
// top left, top right, bottom right, bottom left.
// in most cases, this will end up adding four, but in
// expected cases we may just have three, and in malformed
// cases we may have trapezoids that are segments or points,
// in which ase we will just have two or one returned edge.
func (tr *Trapezoid) DCELEdges() []*dcel.Edge {
	edges := make([]*dcel.Edge, 1)
	i := 0
	edges[i] = dcel.NewEdge()
	edges[i].Origin = tr.Edges[top].Left().Vertex()
	edges[i].Origin.OutEdge = edges[i]
	if !tr.Edges[top].Right().Eq(edges[i].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = tr.Edges[top].Right().Vertex()
		edges[i].Origin.OutEdge = edges[i]
		edges[i-1].Next = edges[i]
		edges[i].Prev = edges[i-1]
	}
	if !tr.Edges[bot].Right().Eq(edges[i].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = tr.Edges[bot].Right().Vertex()
		edges[i].Origin.OutEdge = edges[i]
		edges[i-1].Next = edges[i]
		edges[i].Prev = edges[i-1]
	}
	if !tr.Edges[bot].Left().Eq(edges[i].Origin) &&
		!tr.Edges[bot].Left().Eq(edges[0].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = tr.Edges[bot].Left().Vertex()
		edges[i].Origin.OutEdge = edges[i]
		edges[i-1].Next = edges[i]
		edges[i].Prev = edges[i-1]
	}
	// In the case of a trapezoid which is a point,
	// this will cause the edge to refer to itself by next
	// and prev, which is probably not expected by code
	// which iterates over edges.
	edges[0].Prev = edges[i]
	edges[i].Next = edges[0]
	return edges
}

// Copy returns a trapezoid with identical edges
// and neighbors.
func (tr *Trapezoid) Copy() *Trapezoid {
	tr2 := new(Trapezoid)
	tr2.Edges = tr.Edges
	tr2.Neighbors = tr.Neighbors
	return tr2
}

// HasDefinedPoint returns for a given Trapezoid
// whether or not any of the points on the Trapezoid's
// perimeter match the query point.
// We make an assumption here that there will be no
// edges who have vertices defined on other edges, aka
// that all intersections are represented through
// vertices.
func (tr *Trapezoid) HasDefinedPoint(p dcel.Point) bool {
	for _, e := range tr.Edges {
		for _, p2 := range e {
			if p2 == p {
				return true
			}
		}
	}
	return false
}

func newTrapezoid(sp dcel.Span) *Trapezoid {
	t := new(Trapezoid)
	p1 := *dcel.NewPoint(sp.Min.X(), sp.Max.Y(), sp.Min.Z())
	p2 := *dcel.NewPoint(sp.Max.X(), sp.Min.Y(), sp.Min.Z())
	t.Edges[top] = dcel.FullEdge{sp.Min, p2}
	t.Edges[bot] = dcel.FullEdge{sp.Max, p1}
	t.Edges[left] = dcel.FullEdge{sp.Min, p1}
	t.Edges[right] = dcel.FullEdge{sp.Max, p2}
	return t
}
