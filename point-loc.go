package compgeo

import (
	"errors"
	"sort"

	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree"
)

type LocatesPoints interface {
	PointLocate()
}

type shellNode struct {
	k float64
	v *DCELEdge
}

func (sn shellNode) Key() float64 {
	return sn.k
}

func (sn shellNode) Val() interface{} {
	return sn.v
}

// The real difficulties in Slab Decomposition are all in the
// persistent bst itself, so this is a fairly simple function.
func (dc *DCEL) SlabDecompose(bstType tree.Type) (search.DynamicPersistent, error) {
	t := tree.New(bstType).ToPersistent()
	// Sort points in order of X value
	pts := make([]int, len(dc.Vertices))
	for i := range dc.Vertices {
		pts[i] = i
	}
	if len(dc.Vertices[0]) < 2 {
		// I don't know why someone would want to get the slab decomposition of
		// a structure which has more than two dimensions but there could be
		// applications so we don't reject that idea offhand.
		return nil, errors.New("DCEL's vertices aren't at least two dimensional")
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
			t.Insert(shellNode{v[1], e})
		}
		// // Remove all edges from the PersistentBST connecting to the left
		// // of the point
		for _, e := range rightEdges {
			v2 := e.Twin.Origin
			t.Delete(shellNode{v2[1], e})
		}
	}
	return t, nil
}
