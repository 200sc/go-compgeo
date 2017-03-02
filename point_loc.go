package compgeo

import (
	"errors"
	"sort"
	"strconv"
)

type Point []float64

// Could be called Rect, but Rect implies 2D
type Span struct {
	Start, End Point
}

func (s Span) Length(i int) float64 {
	return s.End[i] - s.Start[i]
}

// If these stay with these names, it needs to
// be in a dcel package
type Edge struct {
	Origin int
	Twin   int
	Face   int
	Next   int
	Prev   int
}

type Face struct {
	Outer, Inner int
}

type DCEL struct {
	Vertices []Point
	// the indices in outEdges map to vertices indices
	// the values in outEdges map to halfEdges indices
	OutEdges  []int
	HalfEdges []Edge
	// The first value in a face is the outside component
	// of the face, the second value is the inside component
	Faces [][2]int
}

func (dc *DCEL) AllEdges(vertex int) []int {
	e1 := dc.OutEdges[vertex]
	edges := make([]int, 1)
	edges[0] = e1
	edge := dc.HalfEdges[dc.HalfEdges[e1].Twin].Next
	for edge != e1 {
		edges = append(edges, edge)
		edge = dc.HalfEdges[dc.HalfEdges[edge].Twin].Next
	}
	return edges
}

// Partition the edges of a vertex by whether they connect
// to a vertex greater or lesser than the given vertex with
// respect to a specific dimension
func (dc *DCEL) PartitionVertexEdges(vertex int, d int) ([]int, []int, error) {
	allEdges := dc.AllEdges(vertex)
	lesser := make([]int, 0)
	greater := make([]int, 0)
	v := dc.Vertices[vertex]
	if len(v) <= d {
		return lesser, greater, errors.New("DCEL's vertex does not support " + strconv.Itoa(d) + " dimensions")
	}
	checkAgainst := v[d]
	for _, i := range allEdges {
		e := dc.HalfEdges[dc.HalfEdges[i].Twin]
		// Potential issue:
		// Will something bad happen if there are multiple
		// elements with the same value in this dimension?
		if dc.Vertices[e.Origin][d] <= checkAgainst {
			lesser = append(lesser, i)
		} else {
			greater = append(greater, i)
		}
	}
	return lesser, greater, nil
}

type Polytope interface {
}

type LocatesPoints interface {
	PointLocate()
}

// The real difficulties in Slab Decomposition are all in the
// persistent bst itself, so this is a fairly simple function.
func (dc *DCEL) SlabDecompose(bstType int) (PersistentBST, error) {
	t := NewPersistentRBTree(bstType)
	// Sort points in order of X value
	pts := make([]int, len(dc.Vertices))
	for i, p := range dc.Vertices {
		pts[i] = i
	}
	if len(dc.Vertices[0]) < 2 {
		// I don't know why someone would want to get the slab decomposition of
		// a structure which has more than two dimensions but there could be
		// applications so we don't reject that idea offhand.
		return nil, errors.New("DCEL's vertices aren't at least  two dimensional")
	}
	// We sort by the 0th dimension here. There is no necessary requirement that
	// the 0th dimension maps to X, but there's also no requirement that slab
	// decomposition uses vertical slabs.
	sort.Slice(pts, func(i, j int) bool {
		return dc.Vertices[pts[i]][0] < dc.Vertices[pts[j]][1]
	})
	// At each point,
	for _, p := range pts {
		v := dc.Vertices[p]
		// Set the BST's instant to the x value of this point
		t.SetInstant(v[0])
		// We don't need to check the returned error here
		// because we already checked this above-- if a DCEL
		// contains points where some points have a different
		// dimension than others that will cause further problems,
		// but this is too expensive to check here.
		leftEdges, rightEdges, _ := dc.PartitionVertexEdges(p, 0)
		// // Add all edges to the PersistentBST connecting to the right
		// // of the point
		for _, e := range leftEdges {
			t.Insert(v[1], e)
		}
		// // Remove all edges from the PersistentBST connecting to the left
		// // of the point
		for _, e := range rightEdges {
			v2 := dc.HalfEdges[dc.HalfEdges[e].Twin].Origin
			t.Delete(dc.Vertices[v2][1], e)
		}
	}
	return t, nil
}
