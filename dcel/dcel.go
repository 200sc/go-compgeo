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

// New returns an empty DCEL with its inner
// fields initialized to empty slices, and a
// zeroth outside face.
func New() *DCEL {
	dc := new(DCEL)
	dc.Vertices = []*Point{}
	dc.OutEdges = []*Edge{}
	dc.HalfEdges = []*Edge{}
	dc.Faces = []*Face{NewFace()}
	return dc
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

// ScanPoints returns the index within
// dc.Vertices where p is, or -1 if it
// does not exist in dc.
func (dc *DCEL) ScanPoints(p *Point) int {
	for i, v := range dc.Vertices {
		if v == p {
			return i
		}
	}
	return -1
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
	//fmt.Println("All edges off of vertex,", dc.Vertices[vertex], "::", allEdges)
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
		if e2.Origin[d] < checkAgainst {
			lesser = append(lesser, e1)
		} else if e2.Origin[d] > checkAgainst {
			greater = append(greater, e1)
		} //else {
		// We completely ignore vertical lines?.
		//}
	}
	return lesser, greater, nil
}
