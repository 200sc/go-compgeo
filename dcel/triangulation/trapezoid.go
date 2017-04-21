package triangulation

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
	Edges     [4]FullEdge
	Neighbors [4]*Trapezoid
	node      *TrapezoidNode
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
func (tr *Trapezoid) HasDefinedPoint(p Point) bool {
	for _, e := range tr.Edges {
		for _, p2 := range e {
			if p2 == p {
				return true
			}
		}
	}
	return false
}
