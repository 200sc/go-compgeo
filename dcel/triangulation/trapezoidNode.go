package triangulation

import (
	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
)

// A TrapezoidNode is a node in a tree structure
// for trapezoid map queries. It is structured so that
// each variety of TrapezoidNode is the same struct, but
// each has a different payload and query function.
type TrapezoidNode struct {
	left, right *TrapezoidNode
	parents     []*TrapezoidNode
	query       func(geom.FullEdge, *TrapezoidNode) []*Trapezoid
	payload     interface{}
}

// DCEL converts the trapezoids in the node search structure
// into a DCEL.
func (tn *TrapezoidNode) DCEL() (*dcel.DCEL, map[*dcel.Face]*dcel.Face) {
	dc := new(dcel.DCEL)
	trs := tn.inOrder()
	// This maps from faces in the output of this algorithm
	// to faces in the input of the TrapezoidMap.
	fMap := make(map[*dcel.Face]*dcel.Face)
	vMap := make(map[geom.Point]*dcel.Edge)
	dc.Faces = make([]*dcel.Face, len(trs))
	// Todo: benchmark if it is faster to initially set this
	// at len(trs) because we know that's a minimum,
	// and then make if checks down the line for whether we append
	// or set.
	dc.Vertices = make([]*dcel.Vertex, 0)
	for i, tr := range trs {
		// Each trapezoid becomes a face
		dc.Faces[i] = dcel.NewFace()
		fMap[dc.Faces[i]] = tr.face
		// each of the up to four edges in a trapezoid
		// has an edge associated with it, the first of which is the face's
		// inner (going with the broken convention of always using inner)
		edges := tr.DCELEdges()
		dc.Faces[i].Inner = edges[0]
		dc.HalfEdges = append(dc.HalfEdges, edges...)
		// each vertex in each trapezoid, if it has not been seen before,
		// is added to the dcel and added to a map connected to it's edge
		// but if it is in the map,
		// it defines that the edge it is mapped to is the twin of the
		// current edge's previous, and vice versa.
		for _, e := range edges {
			e.Face = dc.Faces[i]
			if e2, ok := vMap[e.Origin.Point]; ok {
				e.SetTwin(e2.Prev)
				delete(vMap, e.Origin.Point)
			} else {
				vMap[e.Origin.Point] = e
				dc.Vertices = append(dc.Vertices, e.Origin)
			}
		}
	}
	dc.CorrectTwins()
	return dc, fMap
}

func (tn *TrapezoidNode) inOrder() []*Trapezoid {
	if tn == nil {
		// error, unless this is root,
		// I think
		return []*Trapezoid{}
	}
	if tn.left == nil && tn.right == nil {
		// This is a trapezoid (or should be)
		return []*Trapezoid{tn.payload.(*Trapezoid)}
	}
	trs := tn.left.inOrder()
	return append(trs, tn.right.inOrder()...)
}

// PointLocate returns, from a given complex structure,
// which substructure that point falls into, if any.
func (tn *TrapezoidNode) PointLocate(vs ...float64) (*dcel.Face, error) {
	if len(vs) < 2 {
		return nil, compgeo.InsufficientDimensionsError{}
	}
	// A point query on the structure is equivalent to an
	// edge query where both edges are the same.
	trs := tn.Query(geom.FullEdge{geom.Point{vs[0], vs[1], 0}, geom.Point{vs[0], vs[1], 0}})
	if len(trs) == 0 {
		return nil, nil
	}
	return trs[0].face, nil
}

// Query is shorthand for tn.query(fe, tn)
func (tn *TrapezoidNode) Query(fe geom.FullEdge) []*Trapezoid {
	return tn.query(fe, tn)
}

func (tn *TrapezoidNode) discard(n *TrapezoidNode) {
	for _, p := range tn.parents {
		if p.left == tn {
			p.left = n
		} else {
			p.right = n
		}
	}
	n.parents = tn.parents
	n.parents = []*TrapezoidNode{}
}

func (tn *TrapezoidNode) set(v int, n *TrapezoidNode) {
	switch v {
	case top:
		fallthrough
	case left:
		tn.left = n
	case bot:
		fallthrough
	case right:
		tn.right = n
	}
	n.parents = append(n.parents, tn)
}

// NewRoot returns a root node.
// There is only one root node. We use a root node to avoid
// having "if parent == nil, set this to root" checks in
// our map code.
func NewRoot() *TrapezoidNode {
	return &TrapezoidNode{
		query: rootQuery,
	}
}

func rootQuery(fe geom.FullEdge, n *TrapezoidNode) []*Trapezoid {
	return n.left.Query(fe)
}

// NewX returns an X-Node at point P
func NewX(p geom.D3) *TrapezoidNode {
	return &TrapezoidNode{
		query:   xQuery,
		payload: p,
	}
}

func xQuery(fe geom.FullEdge, n *TrapezoidNode) []*Trapezoid {
	p := n.payload.(geom.Point)
	p2 := p
	p2[1]++
	if geom.IsLeftOf(fe.Left(), p, p2) {
		return n.left.Query(fe)
	}
	return n.right.Query(fe)
}

// NewY returns a Y-Node at edge e
func NewY(e geom.FullEdge) *TrapezoidNode {
	return &TrapezoidNode{
		query:   yQuery,
		payload: e,
	}
}

func yQuery(fe geom.FullEdge, n *TrapezoidNode) []*Trapezoid {
	// This query asks if fe.Left() is above or below
	// yn.FullEdge.
	// If they are colinear, however, we need to check
	// which slope is larger. If fe is larger, we go above,
	// else we go below.
	yn := n.payload.(geom.FullEdge)
	cp := geom.HzCross2D(fe.Left(), yn.Left(), yn.Right())
	if cp > 0 {
		return n.left.Query(fe)
	} else if cp < 0 {
		return n.right.Query(fe)
	}
	// The colinear case
	s1 := fe.Slope()
	s2 := yn.Slope()
	if s1 > s2 {
		return n.left.Query(fe)
	}
	return n.right.Query(fe)
}

// NewTrapNode returns a leaf node holding a trapezoid
func NewTrapNode(tr *Trapezoid) *TrapezoidNode {
	node := &TrapezoidNode{
		payload: tr,
		query:   trapQuery,
	}
	tr.node = node
	return node
}

func trapQuery(fe geom.FullEdge, n *TrapezoidNode) []*Trapezoid {
	tr := n.payload.(*Trapezoid)
	traps := []*Trapezoid{tr}
	r := fe.Right()
	for geom.IsRightOf(r, tr.Edges[right].Left(), tr.Edges[right].Right()) {

		// We perform this check here is it is less expensive
		// than the cross product in the latter case, even
		// though the latter case would suffice to do this.
		if tr.Neighbors[upright] == tr.Neighbors[botright] {
			tr = tr.Neighbors[botright]
		} else {
			// If the edge separating the two
			// trapezoids to the right of tr from one another
			// is a above the query segment, then we also intersect
			// the bottom trapezoid.
			// For this aboveness check we just use the left endpoint
			// of the separating edge, as we know that is within fe's
			// horizontal span.
			if geom.IsAbove(
				tr.Neighbors[upright].Edges[bot].Left(), fe.Left(), fe.Right()) {
				tr = tr.Neighbors[botright]
			} else {
				tr = tr.Neighbors[upright]
			}
		}

		traps = append(traps, tr)
	}
	return traps
}
