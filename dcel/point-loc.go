package dcel

import (
	"errors"
	"fmt"
	"sort"

	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree"
)

type shellNode struct {
	k CompEdge
	v [2]*Face
}

func (sn shellNode) Key() search.Comparable {
	return sn.k
}

func (sn shellNode) Val() interface{} {
	return sn.v
}

type LocatesPoints interface {
	PointLocate(vs ...float64) (*Face, error)
}

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
	fmt.Println(pts)
	for i, p := range pts {
		v := dc.Vertices[p]
		var compX float64
		if i < (len(pts) - 1) {
			v2 := dc.Vertices[pts[i+1]]
			compX = (v[0] + v2[0]) / 2
			compXMap[v[0]] = compX
			fmt.Println("Compx set to", compX, "index", p)
			fmt.Println(dc.Vertices)
		} else {
			// We shouldn't be adding anything at this stage
			compX = -1
		}
		// Set the BST's instant to the x value of this point
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
			compX2 := compXMap[e.Twin.Origin[0]]
			t.Delete(shellNode{CompEdge{e.Twin, compX2}, fs})
			fmt.Println(t)
			//fmt.Println("Removing", e.Twin, "from", v2[1], err)
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
			t.Insert(shellNode{CompEdge{e, compX},
				[2]*Face{e.Face, e.Twin.Face}})
			fmt.Println(t)
			//fmt.Println("Adding", e, "at", v[1])
		}
	}
	return &SlabPointLocator{t, dc.Faces[OUTER_FACE]}, nil
}

// We need to have our keys be CompEdges so
// they are comparable within a certain x range.
type CompEdge struct {
	*Edge
	x float64
}

func (ce CompEdge) Compare(i interface{}) search.CompareResult {
	switch c := i.(type) {
	case CompEdge:
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

type SlabPointLocator struct {
	dp        search.DynamicPersistent
	outerFace *Face
}

func (spl *SlabPointLocator) String() string {
	return fmt.Sprintf("%v", spl.dp)
}

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
	if p.VerticalCompare(e.(CompEdge).Edge) == search.Greater {
		fmt.Println(p, "is above found edge", e)
		return nil, nil
	}

	e2, f2 := tree.SearchUp(p)
	if p.VerticalCompare(e2.(CompEdge).Edge) == search.Less {
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
	// f2 and f1 have both faces in common. We need
	// to find out which face is enclosed in the other.
	// We return the enclosed face.

	// If one is the outer face than we are good
	if f1[0] == spl.outerFace {
		return f1[1], nil
	}
	if f1[1] == spl.outerFace {
		return f1[0], nil
	}
	if f1[0].Encloses(f1[1]) {
		return f1[1], nil
	}
	return f1[0], nil
	// Case VERY unhappy:
	// f2 and f1 have neither face in common, which
	// means something went very wrong and we don't even
	// check for this. This will return in case 1.
}
