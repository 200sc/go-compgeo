package dcel

import (
	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/geom"
)

// A Vertex is a Point which knows its outEdge.
type Vertex struct {
	geom.Point
	OutEdge *Edge
}

// NewVertex returns a Vertex at a given position with no
// outEdge
func NewVertex(x, y, z float64) *Vertex {
	return &Vertex{
		geom.Point{x, y, z},
		nil,
	}
}

// Add adds to the point behind a vertex
func (v *Vertex) Add(d int, f float64) {
	v.Point[d] += f
}

// Mult multiplies the point behind a vertex
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
		err = compgeo.BadDimensionError{}
		return
	}
	allEdges := v.AllEdges()
	checkAgainst := v.Val(d)
	for _, e1 := range allEdges {
		e2 := e1.Twin
		if geom.F64eq(e2.Origin.Val(d), checkAgainst) {
			colinear = append(colinear, e1)
		} else if e2.Origin.Val(d) < checkAgainst {
			lesser = append(lesser, e1)
		} else if e2.Origin.Val(d) > checkAgainst {
			greater = append(greater, e1)
		}
	}
	return
}

// PointToVertex converts a point into a vertex
func PointToVertex(dp geom.D3) *Vertex {
	return NewVertex(dp.Val(0), dp.Val(1), dp.Val(2))
}
