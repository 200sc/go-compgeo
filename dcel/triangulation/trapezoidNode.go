package triangulation

import (
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
	query       func(FullEdge, *TrapezoidNode) []*Trapezoid
	payload     interface{}
}

// Query is shorthand for n.query(fe, n)
func (tn *TrapezoidNode) Query(fe FullEdge) []*Trapezoid {
	return tn.query(fe, tn)
}

func (tn *TrapezoidNode) Discard(n *TrapezoidNode) {
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

func (tn *TrapezoidNode) Set(v int, n *TrapezoidNode) {
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

func rootQuery(fe FullEdge, n *TrapezoidNode) []*Trapezoid {
	return n.left.Query(fe)
}

// NewX returns an X-Node at point P
func NewX(p Point) *TrapezoidNode {
	return &TrapezoidNode{
		query:   xQuery,
		payload: p,
	}
}

func xQuery(fe FullEdge, n *TrapezoidNode) []*Trapezoid {
	p := n.payload.(dcel.Point)
	p2 := p
	p2[1]++
	if IsLeftOf(fe.Left(), p, p2) {
		return n.left.Query(fe)
	}
	return n.right.Query(fe)
}

// NewY returns a Y-Node at edge e
func NewY(e FullEdge) *TrapezoidNode {
	return &TrapezoidNode{
		query:   yQuery,
		payload: e,
	}
}

func yQuery(fe FullEdge, n *TrapezoidNode) []*Trapezoid {
	// This query asks if fe.Left() is above or below
	// yn.FullEdge.
	// If they are colinear, however, we need to check
	// which slope is larger. If fe is larger, we go above,
	// else we go below.
	yn := n.payload.(FullEdge)
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

func trapQuery(fe FullEdge, n *TrapezoidNode) []*Trapezoid {
	tr := n.payload.(*Trapezoid)
	traps := []*Trapezoid{tr}
	r := fe.Right()
	for IsRightOf(r, tr.Edges[right].Left(), tr.Edges[right].Right()) {

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
			if IsAbove(
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
