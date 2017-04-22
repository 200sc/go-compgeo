package slab

import (
	"fmt"
	"image/color"

	"bitbucket.org/oakmoundstudio/oak/physics"

	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/visualize"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree"
)

type faces struct {
	f1, f2 *dcel.Face
}

func (fs faces) Equals(e search.Equalable) bool {
	switch fs2 := e.(type) {
	case faces:
		return fs2.f1 == fs.f1 && fs2.f2 == fs.f2
	}
	return false
}

type shellNode struct {
	k compEdge
	v search.Equalable
}

func (sn shellNode) Key() search.Comparable {
	return sn.k
}

func (sn shellNode) Val() search.Equalable {
	return sn.v
}

// Decompose is based on Dobkin and Lipton's work into
// point location.
// The real difficulties in Slab Decomposition are all in the
// persistent bst itself, so this is a fairly simple function.
func Decompose(dc *dcel.DCEL, bstType tree.Type) (dcel.LocatesPoints, error) {
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
	compXMap := make(map[float64]float64)
	i := 1
OUTER:
	for i < len(pts) {
		j := i - 1
		prevX := dc.Vertices[pts[j]].X()
		thisX := dc.Vertices[pts[i]].X()
		for geom.F64eq(thisX, prevX) {
			i++
			if i == len(pts) {
				break OUTER
			}
			thisX = dc.Vertices[pts[i]].X()
		}
		compXMap[geom.ToFixed(prevX, 5)] = (prevX + thisX) / 2
		i++
	}
	compXMap[dc.Vertices[pts[len(pts)-1]].X()] = -1

	fmt.Println(dc.Vertices)
	fmt.Println(pts)
	fmt.Println(compXMap)
	fmt.Println("START CONSTRUCTION")
	i = 0
	for i < len(pts) {
		p := pts[i]
		v := dc.Vertices[p]
		// Set the BST's instant to the x value of this point
		fmt.Println("Setting Instant to", v.X())
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
		if visualize.VisualCh != nil {
			visualize.HighlightColor = color.RGBA{255, 0, 0, 255}
		}
		for _, e := range le {
			fmt.Println("Removing", e.Twin, compXMap[e.Twin.Origin.X()], compXMap[geom.ToFixed(e.Twin.Origin.X(), 5)])
			err := ct.Delete(shellNode{compEdge{e.Twin, compXMap[geom.ToFixed(e.Twin.Origin.X(), 5)]}, search.Nil{}})
			fmt.Println("Remove result", err)
			fmt.Println(ct)
		}
		// Add all edges to the PersistentBST connecting to the right
		// of the point
		if visualize.VisualCh != nil {
			visualize.HighlightColor = color.RGBA{0, 255, 0, 255}
		}
		for _, e := range re {
			// We always want the half edge that points to the right,
			// and between the two faces this edge is on we want the
			// face which is LOWER. This is because we ulimately point
			// locate to the edge above the query point. Returning an
			// edge for a query represents that the query is below
			// the edge,
			fmt.Println("Adding", e, "at", geom.ToFixed(v.X(), 5))
			ct.Insert(shellNode{compEdge{e, compXMap[geom.ToFixed(v.X(), 5)]},
				faces{e.Face, e.Twin.Face}})
			fmt.Println(ct)
		}

		i++
	}
	fmt.Println("END CONSTRUCTION")
	if visualize.VisualCh != nil {
		visualize.HighlightColor = color.RGBA{255, 255, 255, 255}
	}
	return &PointLocator{t, dc.Faces[dcel.OUTER_FACE]}, nil
}

// We need to have our keys be CompEdges so
// they are comparable within a certain x range.
type compEdge struct {
	*dcel.Edge
	x float64
}

func (ce compEdge) Compare(i interface{}) search.CompareResult {
	switch c := i.(type) {
	case compEdge:
		if visualize.VisualCh != nil {
			visualize.DrawLine(ce.Edge.Origin, ce.Edge.Twin.Origin)
			visualize.DrawLine(c.Edge.Origin, c.Edge.Twin.Origin)
		}
		fmt.Println("Comparing", ce, c)
		if ce.Edge == c.Edge {
			fmt.Println("Equal1!")
			return search.Equal
		}

		if geom.F64eq(ce.X(), c.X()) && geom.F64eq(ce.Y(), c.Y()) &&
			geom.F64eq(ce.Twin.X(), c.Twin.X()) && geom.F64eq(ce.Twin.Y(), c.Twin.Y()) {
			fmt.Println("Equal2!")
			return search.Equal
		}
		compX := ce.x
		compXbackup := c.x
		if c.x > compX {
			compX = c.x
			compXbackup = ce.x
		}
		tryBackup := false
		p1, err := ce.PointAt(0, compX)
		if err != nil {
			fmt.Println("compX", compX, "not on point ", ce)
			tryBackup = true
		}
		p2, err := c.PointAt(0, compX)
		if err != nil {
			fmt.Println("compX", compX, "not on point ", c)
			tryBackup = true
		}
		if tryBackup {
			compX = compXbackup
			p1, err = ce.PointAt(0, compX)
			if err != nil {
				fmt.Println("Backup", compX, "not on point ", ce)
			}
			p2, err = c.PointAt(0, compX)
			if err != nil {
				fmt.Println("Backup", compX, "not on point ", c)
			}
		}
		if p1[1] < p2[1] {
			fmt.Println("Less!")
			return search.Less
		}
		fmt.Println("Greater!")
		return search.Greater
	}
	return ce.Edge.Compare(i)
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
	fmt.Println("Searching on tree")
	e, f := tree.SearchDown(p, 0)
	if e == nil {
		fmt.Println("Location on empty tree")
		return nil, nil
	}
	fmt.Println("Edge found", e)
	if geom.VerticalCompare(p, e.(compEdge)) == search.Greater {
		fmt.Println(p, "is above found edge", e)
		return nil, nil
	}

	e2, f2 := tree.SearchUp(p)
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
			if visualize.VisualCh != nil {
				ps := f5.Vertices()
				physVerts := make([]physics.Vector, len(ps))
				for i, v := range ps {
					physVerts[i] = physics.NewVector(v.X(), v.Y())
				}
				visualize.DrawPoly(physVerts)
			}
			if f5.Contains(p) {
				fmt.Println("P was contained")
				return f5, nil
			}
		}
	}

	return nil, nil
}
