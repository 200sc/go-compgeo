// slab implements point location by splitting a dcel into vertical slabs

package slab

import (
	"fmt"

	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree"
)

// Decompose is based on Dobkin and Lipton's work into
// point location.
// The real difficulties in Slab Decomposition are all in the
// persistent bst itself, so this is a fairly simple function.
func Decompose(dc *dcel.DCEL, bstType tree.Type) (pointLoc.LocatesPoints, error) {
	if dc == nil || len(dc.Vertices) < 3 {
		return nil, compgeo.BadDCELError{}
	}
	if dc.Vertices[0].D() < 2 {
		// I don't know why someone would want to get the slab decomposition of
		// a structure which has more than two dimensions but there could be
		// applications so we don't reject that idea offhand.
		return nil, compgeo.BadDimensionError{}
	}
	t := tree.New(bstType).ToPersistent()
	pts := dc.VerticesSorted(0)

	// For each edge, we need to know which face lies beneath it in its two
	// faces.
	faceEdgeMap := make(map[*dcel.Edge]*dcel.Face)
	for _, f := range dc.Faces {
		// walk each face
		e := f.Outer
		if e != nil {
			if e.Origin.X() < e.Twin.Origin.X() {
				faceEdgeMap[e] = f
			}
			for e = e.Next; e != f.Outer; e = e.Next {
				// This edge points right, in an outer face.
				// Then this face lies beneath e.
				if e.Origin.X() < e.Twin.Origin.X() {
					faceEdgeMap[e] = f
				}
			}
		}
		e = f.Inner
		if e != nil {
			if e.Origin.X() > e.Twin.Origin.X() {
				faceEdgeMap[e] = f
			}
			for e = e.Next; e != f.Inner; e = e.Next {
				if e.Origin.X() > e.Twin.Origin.X() {
					faceEdgeMap[e] = f
				}
			}
		}
	}

	i := 0
	for i < len(pts) {
		p := pts[i]
		v := dc.Vertices[p]
		// Set the BST's instant to the x value of this point
		t.SetInstant(v.X())
		ct := t.ThisInstant()

		// Aggregate all points at this x value so we do not
		// attempt to add edges to a tree which contains edges
		// point to the left of v[0]
		vs := []*dcel.Vertex{v}
		for (i+1) < len(pts) && geom.F64eq(dc.Vertices[pts[i+1]].X(), v.X()) {
			i++
			p = pts[i]
			vs = append(vs, dc.Vertices[p])
		}

		le := []*dcel.Edge{}
		re := []*dcel.Edge{}

		for _, v := range vs {
			// We don't need to check the returned error here
			// because we already checked this above-- if a DCEL
			// contains points where some points have a different
			// dimension than others that will cause further problems,
			// but this is too expensive to check here.
			leftEdges, rightEdges, _, _ := v.PartitionEdges(0)
			le = append(le, leftEdges...)
			re = append(re, rightEdges...)
		}
		// Remove all edges from the PersistentBST connecting to the left
		// of the points
		for _, e := range le {
			err := ct.Delete(shellNode{compEdge{e.Twin}, search.Nil{}})
			if err != nil {
				fmt.Println(err, e.Twin)
			}
		}
		// Add all edges to the PersistentBST connecting to the right
		// of the point
		for _, e := range re {
			// We always want the half edge that points to the right,
			// and between the two faces this edge is on we want the
			// face which is LOWER. This is because we ulimately point
			// locate to the edge above the query point. Returning an
			// edge for a query represents that the query is below
			// the edge,
			ct.Insert(shellNode{compEdge{e}, face{faceEdgeMap[e]}})
		}

		i++
	}
	return &PointLocator{t, dc.Faces[dcel.OUTER_FACE]}, nil
}

// PointLocator is a construct that uses slab
// decomposition for point location.
type PointLocator struct {
	dp        search.DynamicPersistent
	outerFace *dcel.Face
}

func (spl *PointLocator) String() string {
	return fmt.Sprintf("%v", spl.dp)
}

// PointLocate returns which face within this SlabPointLocator
// the query point lands, within two dimensions.
func (spl *PointLocator) PointLocate(vs ...float64) (*dcel.Face, error) {
	if len(vs) < 2 {
		return nil, compgeo.InsufficientDimensionsError{}
	}
	tree := spl.dp.AtInstant(vs[0])
	p := geom.Point{vs[0], vs[1], 0}

	_, f2 := tree.SearchUp(p, 0)

	return f2.(face).Face, nil
}
