package dcel

import (
	"errors"
	"fmt"
	"sort"

	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree"
)

type shellNode struct {
	k compEdge
	v [2]*Face
}

func (sn shellNode) Key() search.Comparable {
	return sn.k
}

func (sn shellNode) Val() interface{} {
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
		return dc.Vertices[pts[i]][0] < dc.Vertices[pts[j]][0]
	})
	compXMap := make(map[float64]float64)
	i := 1
OUTER:
	for i < len(pts) {
		j := i - 1
		prevX := dc.Vertices[pts[j]][0]
		thisX := dc.Vertices[pts[i]][0]
		for thisX == prevX {
			i++
			if i == len(pts) {
				break OUTER
			}
			thisX = dc.Vertices[pts[i]][0]
		}
		compXMap[prevX] = (prevX + thisX) / 2
		i++
	}
	compXMap[dc.Vertices[pts[len(pts)-1]][0]] = -1

	fmt.Println(dc.Vertices)
	fmt.Println(pts)
	fmt.Println(compXMap)
	fmt.Println("START CONSTRUCTION")
	for _, p := range pts {
		v := dc.Vertices[p]
		// Set the BST's instant to the x value of this point
		fmt.Println("Setting Instant to", v[0])
		t.SetInstant(v[0])

		// We don't need to check the returned error here
		// because we already checked this above-- if a DCEL
		// contains points where some points have a different
		// dimension than others that will cause further problems,
		// but this is too expensive to check here.
		leftEdges, rightEdges, _ := dc.PartitionVertexEdges(p, 0)
		// Remove all edges from the PersistentBST connecting to the left
		// of the point
		for _, e := range leftEdges {
			var fs [2]*Face
			fmt.Println("Removing", e.Twin)
			err := t.Delete(shellNode{compEdge{e.Twin, compXMap[e.Twin.Origin[0]]}, fs})
			fmt.Println("Removed", err)
			fmt.Println(t)

		}
		// Add all edges to the PersistentBST connecting to the right
		// of the point
		for _, e := range rightEdges {
			// We always want the half edge that points to the right,
			// and between the two faces this edge is on we want the
			// face which is LOWER. This is because we ulimately point
			// locate to the edge above the query point. Returning an
			// edge for a query represents that the query is below
			// the edge,
			fmt.Println("Adding", e, "at", v[0])
			t.Insert(shellNode{compEdge{e, compXMap[v[0]]},
				[2]*Face{e.Face, e.Twin.Face}})
			fmt.Println(t)
		}
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

		if p1[1] == p2[1] {
			return search.Equal
		} else if p1[1] < p2[1] {
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
		return nil, errors.New("Slab point location only supports 2 dimensions")
	}
	fmt.Println("Querying", vs)
	tree := spl.dp.AtInstant(vs[0])
	fmt.Println("Tree found:")
	fmt.Println(tree)
	p := Point{vs[0], vs[1], 0}
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
	f1 := f.([2]*Face)
	f3 := f2.([2]*Face)
	if f1[0] != f3[0] && f1[0] != f3[1] {
		return f1[1], nil
	} else if f1[1] != f3[0] && f1[1] != f3[1] {
		return f1[0], nil
	}
	// Case unhappy:
	// f2 and f1 have both faces in common.
	// We then do PIP on each face, and return
	// whichever is true, if either.
	fmt.Println("Checking if face contains", p)
	if f1[0] != spl.outerFace && f1[0].Contains(p) {
		fmt.Println("P was contained")
		return f1[0], nil
	}
	fmt.Println("Checking if other face contains", p)
	if f1[1] != spl.outerFace && f1[1].Contains(p) {
		fmt.Println("P was contained")
		return f1[1], nil
	}
	return nil, nil
	// Case VERY unhappy:
	// f2 and f1 have neither face in common, which
	// means something went very wrong and we don't even
	// check for this. This will return in case 1.
}
