package dcel

// A Vertex is a Point which knows its outEdge.
type Vertex struct {
	Point
	OutEdge *Edge
}

// NewVertex returns a Vertex at a given position with no
// outEdge
func NewVertex(x, y, z float64) *Vertex {
	return &Vertex{
		Point{x, y, z},
		nil,
	}
}

func (v *Vertex) Set(d int, f float64) {
	v.Point[d] = f
}

func (v *Vertex) Add(d int, f float64) {
	v.Point[d] += f
}

func (v *Vertex) Mult(d int, f float64) {
	v.Point[d] *= f
}

// AllEdges iterates through the edges surrounding
// a vertex and returns them all.
func (v *Vertex) AllEdges() []*Edge {
	return v.OutEdge.AllEdges()
}

// PartitionEdges splits the edges around a vertex and
// returns those whose endpoints are lesser, equal to, and greater
// than the given vertex in dimension d.
func (v *Vertex) PartitionEdges(d int) (lesser []*Edge,
	greater []*Edge, colinear []*Edge, err error) {

	if len(v.Point) <= d {
		err = BadDimensionError{}
		return
	}
	allEdges := v.AllEdges()
	checkAgainst := v.Val(d)
	for _, e1 := range allEdges {
		e2 := e1.Twin
		// Potential issue:
		// Will something bad happen if there are multiple
		// elements with the same value in this dimension?
		// Answer: Yes yes yes
		if f64eq(e2.Origin.Val(d), checkAgainst) {
			colinear = append(colinear, e1)
		} else if e2.Origin.Val(d) < checkAgainst {
			lesser = append(lesser, e1)
		} else if e2.Origin.Val(d) > checkAgainst {
			greater = append(greater, e1)
		}
	}
	return
}
