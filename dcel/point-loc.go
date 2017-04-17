package dcel

import (
	"fmt"
	"sort"

	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree"
)

type faces struct {
	f1, f2 *Face
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

// LocatesPoints is an interface to represent point location
// queries.
type LocatesPoints interface {
	PointLocate(vs ...float64) (*Face, error)
}

// SlabDecompose is based on Dobkin and Lipton's work into
// point location.
// The real difficulties in Slab Decomposition are all in the
// persistent bst itself, so this is a fairly simple function.
func (dc *DCEL) SlabDecompose(bstType tree.Type) (LocatesPoints, error) {
	if dc == nil || len(dc.Vertices) < 3 {
		return nil, BadDCELError{}
	}
	t := tree.New(bstType).ToPersistent()
	// Sort points in order of X value
	pts := make([]int, len(dc.Vertices))
	for i := range dc.Vertices {
		pts[i] = i
	}
	if dc.Vertices[0].D() < 2 {
		// I don't know why someone would want to get the slab decomposition of
		// a structure which has more than two dimensions but there could be
		// applications so we don't reject that idea offhand.
		return nil, BadDimensionError{}
	}
	// We sort by the 0th dimension here. There is no necessary requirement that
	// the 0th dimension maps to X, but there's also no requirement that slab
	// decomposition uses vertical slabs.
	sort.Slice(pts, func(i, j int) bool {
		return dc.Vertices[pts[i]].X() < dc.Vertices[pts[j]].X()
	})
	compXMap := make(map[float64]float64)
	i := 1
OUTER:
	for i < len(pts) {
		j := i - 1
		prevX := dc.Vertices[pts[j]].X()
		thisX := dc.Vertices[pts[i]].X()
		for thisX == prevX {
			i++
			if i == len(pts) {
				break OUTER
			}
			thisX = dc.Vertices[pts[i]].X()
		}
		compXMap[prevX] = (prevX + thisX) / 2
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

		// Aggregate all points at this x value so we do not
		// attempt to add edges to a tree which contains edges
		// point to the left of v[0]
		vs := []*Vertex{v}
		for (i+1) < len(pts) && dc.Vertices[pts[i+1]].X() == v.X() {
			i++
			p = pts[i]
			vs = append(vs, dc.Vertices[p])
		}

		le := []*Edge{}
		re := []*Edge{}

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
		for _, e := range le {
			fmt.Println("Removing", e.Twin)
			err := t.Delete(shellNode{compEdge{e.Twin, compXMap[e.Twin.Origin.X()]}, search.Nil{}})
			fmt.Println("Remove result", err)
			fmt.Println(t)

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
			fmt.Println("Adding", e, "at", v.X())
			t.Insert(shellNode{compEdge{e, compXMap[v.X()]},
				faces{e.Face, e.Twin.Face}})
			fmt.Println(t)
		}

		i++
	}
	fmt.Println("END CONSTRUCTION")
	return &SlabPointLocator{t, dc.Faces[OUTER_FACE]}, nil
}

// We need to have our keys be CompEdges so
// they are comparable within a certain x range.
type compEdge struct {
	*Edge
	x float64
}

func (ce compEdge) Compare(i interface{}) search.CompareResult {
	switch c := i.(type) {
	case compEdge:
		fmt.Println("Comparing", ce, c)
		if ce.Edge == c.Edge {
			return search.Equal
		}
		if ce.X() == c.X() && ce.Y() == c.Y() &&
			ce.Twin.X() == c.Twin.X() && ce.Twin.Y() == c.Twin.Y() {
			return search.Equal
		}
		compX := ce.x
		if c.x > compX {
			compX = c.x
		}
		p1, err := ce.PointAt(0, compX)
		if err != nil {
			fmt.Println("compX", compX, "not on point ", ce)
		}
		p2, err := c.PointAt(0, compX)
		if err != nil {
			fmt.Println("compX", compX, "not on point ", c)
		}
		if p1[1] < p2[1] {
			return search.Less
		}
		return search.Greater
	}
	return ce.Edge.Compare(i)
}

// SlabPointLocator is a construct that uses slab
// decomposition for point location.
type SlabPointLocator struct {
	dp        search.DynamicPersistent
	outerFace *Face
}

func (spl *SlabPointLocator) String() string {
	return fmt.Sprintf("%v", spl.dp)
}

// PointLocate returns which face within this SlabPointLocator
// the query point lands, within two dimensions.
func (spl *SlabPointLocator) PointLocate(vs ...float64) (*Face, error) {
	if len(vs) < 2 {
		return nil, InsufficientDimensionsError{}
	}
	fmt.Println("Querying", vs)
	tree := spl.dp.AtInstant(vs[0])
	fmt.Println("Tree found:")
	fmt.Println(tree)
	p := Point{vs[0], vs[1], 0}
	fmt.Println("Searching on tree")
	e, f := tree.SearchDown(p)
	if e == nil {
		fmt.Println("Location on empty tree")
		return nil, nil
	}
	fmt.Println("Edge found", e)
	if p.VerticalCompare(e.(compEdge).Edge) == search.Greater {
		fmt.Println(p, "is above found edge", e)
		return nil, nil
	}

	e2, f2 := tree.SearchUp(p)
	if p.VerticalCompare(e2.(compEdge).Edge) == search.Less {
		fmt.Println(p, "is below edge", e2)
		return nil, nil
	}
	// Case Happy:
	// f2 and f1 have one face in common. Return it.
	f1 := f.(faces)
	f3 := f2.(faces)
	if f1.f1 != f3.f1 && f1.f1 != f3.f2 {
		return f1.f2, nil
	} else if f1.f2 != f3.f1 && f1.f2 != f3.f2 {
		return f1.f1, nil
	}
	// Case unhappy:
	// f2 and f1 have both faces in common.
	// We then do PIP on each face, and return
	// whichever is true, if either.
	f4 := f1.f2
	f5 := f1.f1

	fmt.Println("Checking if face contains", p)
	if f5 != spl.outerFace && f5.Contains(p) {
		fmt.Println("P was contained")
		return f5, nil
	}
	fmt.Println("Checking if other face contains", p)
	if f4 != spl.outerFace && f4.Contains(p) {
		fmt.Println("P was contained")
		return f4, nil
	}
	return nil, nil
	// Case VERY unhappy:
	// f2 and f1 have neither face in common, which
	// means something went very wrong and we don't even
	// check for this. This will return in case 1.
}
