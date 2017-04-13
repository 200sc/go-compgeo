package dcel

import (
	"errors"
	"math"
	"strconv"
)

// A DCEL is a structure representin arbitrary plane
// divisions and 3d polytopes. Its values are relatively
// self-explanatory but constructing it is significantly
// harder.
type DCEL struct {
	Vertices []*Point
	// outEdges[0] is the (an) edge in HalfEdges whose
	// orgin is Vertices[0]
	OutEdges  []*Edge
	HalfEdges []*Edge
	// The first value in a face is the outside component
	// of the face, the second value is the inside component
	Faces []*Face
}

// MaxX returns the Maximum of all X values
func (dc *DCEL) MaxX() float64 {
	return dc.Max(0)
}

// MaxY returns the Maximum of all Y values
func (dc *DCEL) MaxY() float64 {
	return dc.Max(1)
}

// MaxZ returns the Maximum of all Z values
func (dc *DCEL) MaxZ() float64 {
	return dc.Max(2)
}

// Max functions iterate through vertices
// to find the maximum value along a given axis
// in the DCEL
func (dc *DCEL) Max(i int) (x float64) {
	for _, p := range dc.Vertices {
		if p[i] > x {
			x = p[i]
		}
	}
	return x
}

// MinX returns the Minimum of all X values
func (dc *DCEL) MinX() float64 {
	return dc.Min(0)
}

// MinY returns the Minimum of all Y values
func (dc *DCEL) MinY() float64 {
	return dc.Min(1)
}

// MinZ returns the Minimum of all Z values
func (dc *DCEL) MinZ() float64 {
	return dc.Min(2)
}

// Min functions iterate through vertices
// to find the maximum value along a given axis
// in the DCEL
func (dc *DCEL) Min(i int) (x float64) {
	x = math.Inf(1)
	for _, p := range dc.Vertices {
		if p[i] < x {
			x = p[i]
		}
	}
	return x
}

// AllEdges iterates through the edges surrounding
// a vertex and returns them all.
func (dc *DCEL) AllEdges(vertex int) []*Edge {
	e1 := dc.OutEdges[vertex]
	edges := make([]*Edge, 1)
	edges[0] = e1
	edge := e1.Twin.Next
	for edge != e1 {
		edges = append(edges, edge)
		edge = edge.Twin.Next
	}
	return edges
}

// PartitionVertexEdges partitions the edges of a vertex by
// whether they connect to a vertex greater or lesser than the
// given vertex with respect to a specific dimension
func (dc *DCEL) PartitionVertexEdges(vertex int, d int) ([]*Edge, []*Edge, error) {
	allEdges := dc.AllEdges(vertex)
	lesser := make([]*Edge, 0)
	greater := make([]*Edge, 0)
	v := dc.Vertices[vertex]
	if len(v) <= d {
		return lesser, greater, errors.New("DCEL's vertex does not support " + strconv.Itoa(d) + " dimensions")
	}
	checkAgainst := v[d]
	for _, e1 := range allEdges {
		e2 := e1.Twin
		// Potential issue:
		// Will something bad happen if there are multiple
		// elements with the same value in this dimension?
		if e2.Origin[d] <= checkAgainst {
			lesser = append(lesser, e1)
		} else {
			greater = append(greater, e1)
		}
	}
	return lesser, greater, nil
}
