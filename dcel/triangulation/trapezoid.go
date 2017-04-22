package triangulation

import (
	"bitbucket.org/oakmoundstudio/oak/physics"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/printutil"
)

// These constants refer to indices
// within trapezoids' Edges
const (
	top = iota
	bot
)
const (
	left = iota
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
	top         [2]float64 // y values
	bot         [2]float64 // y values
	left, right float64    // x values
	Neighbors   [4]*Trapezoid
	node        *TrapezoidNode
	faces       [2]*dcel.Face
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
	edges[i].Origin = dcel.PointToVertex(geom.NewPoint(tr.left, tr.top[left], 0))
	edges[i].Origin.OutEdge = edges[i]
	p := geom.NewPoint(tr.right, tr.top[right], 0)
	if !p.Eq(edges[i].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = dcel.PointToVertex(p)
		edges[i].Origin.OutEdge = edges[i]
		edges[i-1].Next = edges[i]
		edges[i].Prev = edges[i-1]
	}
	p = geom.NewPoint(tr.right, tr.bot[right], 0)
	if !p.Eq(edges[i].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = dcel.PointToVertex(p)
		edges[i].Origin.OutEdge = edges[i]
		edges[i-1].Next = edges[i]
		edges[i].Prev = edges[i-1]
	}
	p = geom.NewPoint(tr.left, tr.bot[left], 0)
	if !p.Eq(edges[i].Origin) &&
		!p.Eq(edges[0].Origin) {
		i++
		edges = append(edges, dcel.NewEdge())
		edges[i].Origin = dcel.PointToVertex(p)
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

// Rights is shorthand for setting both of
// tr's right neighbors to the same value
func (tr *Trapezoid) Rights(tr2 *Trapezoid) {
	tr.Neighbors[upright] = tr2
	tr.Neighbors[botright] = tr2
}

// Lefts is shorthand for setting both of
// tr's left neighbors to the same value.
func (tr *Trapezoid) Lefts(tr2 *Trapezoid) {
	tr.Neighbors[upleft] = tr2
	tr.Neighbors[botleft] = tr2
}

func (tr *Trapezoid) setRight(x float64) {
	tr.right = x
}

func (tr *Trapezoid) setLeft(x float64) {
	tr.left = x
}

// Copy returns a trapezoid with identical edges
// and neighbors.
func (tr *Trapezoid) Copy() *Trapezoid {
	tr2 := new(Trapezoid)
	tr2.top = tr.top
	tr2.bot = tr.bot
	tr2.left = tr.left
	tr2.right = tr.right
	tr2.Neighbors = tr.Neighbors
	return tr2
}

// AsPoints converts a trapezoid's internal values
// into four points.
func (tr *Trapezoid) AsPoints() []geom.D2 {
	ds := make([]geom.D2, 4)
	ds[0] = geom.NewPoint(tr.left, tr.top[left], 0)
	ds[1] = geom.NewPoint(tr.right, tr.top[right], 0)
	ds[2] = geom.NewPoint(tr.right, tr.bot[right], 0)
	ds[3] = geom.NewPoint(tr.left, tr.bot[left], 0)
	return ds
}

// BotEdge returns a translation of tr's values to
// tr's bottom edge as a FullEdge
func (tr *Trapezoid) BotEdge() geom.FullEdge {
	return geom.FullEdge{
		geom.NewPoint(tr.right, tr.bot[right], 0),
		geom.NewPoint(tr.left, tr.bot[left], 0),
	}
}

// TopEdge acts as BotEdge for tr's top
func (tr *Trapezoid) TopEdge() geom.FullEdge {
	return geom.FullEdge{
		geom.NewPoint(tr.left, tr.top[left], 0),
		geom.NewPoint(tr.right, tr.top[right], 0),
	}
}

// HasDefinedPoint returns for a given Trapezoid
// whether or not any of the points on the Trapezoid's
// perimeter match the query point.
// We make an assumption here that there will be no
// edges who have vertices defined on other edges, aka
// that all intersections are represented through
// vertices.
func (tr *Trapezoid) HasDefinedPoint(p geom.D3) bool {
	for _, p2 := range tr.AsPoints() {
		if p2.Eq(p) {
			return true
		}
	}
	return false
}

func (tr *Trapezoid) String() string {
	s := ""
	for _, p := range tr.AsPoints() {
		s += "("
		s += printutil.Stringf64(p.X(), p.Y())
		s += ")"
	}
	return s
}

func newTrapezoid(sp geom.Span) *Trapezoid {
	t := new(Trapezoid)
	min := sp.At(geom.SPAN_MIN).(geom.Point)
	max := sp.At(geom.SPAN_MAX).(geom.Point)
	t.top[left] = max.Y()
	t.top[right] = max.Y()
	t.bot[left] = min.Y()
	t.bot[right] = min.Y()
	t.left = min.X()
	t.right = max.X()
	t.Neighbors = [4]*Trapezoid{}
	return t
}

func (tr *Trapezoid) toPhysics() []physics.Vector {
	vs := make([]physics.Vector, 4)
	for i, p := range tr.AsPoints() {
		vs[i] = physics.NewVector(p.X(), p.Y())
	}
	return vs
}
