// slab implements point location by splitting a dcel into vertical slabs

package slab

import (
	"fmt"

	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc"
	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
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

	i := 0
	for i < len(pts) {
		p := pts[i]
		v := dc.Vertices[p]
		// Set the BST's instant to the x value of this point
		visualize.HighlightColor = visualize.CheckLineColor
		visualize.DrawVerticalLine(v)
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
		fmt.Println("Left Edges", le)
		fmt.Println("Right Edges", re)
		// Remove all edges from the PersistentBST connecting to the left
		// of the points
		visualize.HighlightColor = visualize.RemoveColor
		for _, e := range le {
			fmt.Println("Removing", e.Twin)
			err := ct.Delete(shellNode{compEdge{e.Twin}, search.Nil{}})
			fmt.Println("Remove result", err)
			fmt.Println(ct)
		}
		// Add all edges to the PersistentBST connecting to the right
		// of the point
		visualize.HighlightColor = visualize.AddColor
		for _, e := range re {
			// We always want the half edge that points to the right,
			// and between the two faces this edge is on we want the
			// face which is LOWER. This is because we ulimately point
			// locate to the edge above the query point. Returning an
			// edge for a query represents that the query is below
			// the edge,
			fmt.Println("Adding", e)
			ct.Insert(shellNode{compEdge{e}, faces{e.Face, e.Twin.Face}})
			fmt.Println(ct)
		}

		i++
	}
	visualize.HighlightColor = visualize.CheckLineColor
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
	fmt.Println("Querying", vs)
	tree := spl.dp.AtInstant(vs[0])
	fmt.Println("Tree found:")
	fmt.Println(tree)
	p := geom.Point{vs[0], vs[1], 0}

	e, f := tree.SearchDown(p, 0)
	if e == nil {
		fmt.Println("Location on empty tree")
		return nil, nil
	}
	e2, f2 := tree.SearchUp(p, 0)
	fmt.Println("Edges found", e, e2)
	if geom.VerticalCompare(p, e.(compEdge)) == search.Greater {
		fmt.Println(p, "is above edge", e)
		return nil, nil
	}

	if geom.VerticalCompare(p, e2.(compEdge)) == search.Less {
		fmt.Println(p, "is below edge", e2)
		return nil, nil
	}

	// We then do PIP on each face, and return
	// whichever is true, if any.
	f3 := f.(faces)
	f4 := f2.(faces)
	faces := []*dcel.Face{f3.f1, f3.f2, f4.f1, f4.f2}

	for _, f5 := range faces {
		if f5 != spl.outerFace {
			fmt.Println("Checking if face contains", p)
			visualize.HighlightColor = visualize.CheckFaceColor
			visualize.DrawFace(f5)
			if f5.Contains(p) {
				fmt.Println("P was contained")
				return f5, nil
			}
		}
	}

	return nil, nil
}
