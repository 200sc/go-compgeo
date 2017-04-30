package trapezoid

import "github.com/200sc/go-compgeo/geom"

// NewTrapNode returns a leaf node holding a trapezoid
func NewTrapNode(tr *Trapezoid) *Node {
	node := &Node{
		payload: tr,
		query:   trapQuery,
	}
	tr.node = node
	return node
}

func trapQuery(fe geom.FullEdge, n *Node) []*Trapezoid {
	tr := n.payload.(*Trapezoid)
	traps := []*Trapezoid{tr}
	r := fe.Right()
	for tr != nil && r.X() > tr.right {
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
				tr.Neighbors[upright].BotEdge().Left(), fe.Left(), fe.Right()) {
				tr = tr.Neighbors[botright]
			} else {
				tr = tr.Neighbors[upright]
			}
		}
		if tr != nil {
			traps = append(traps, tr)
		}
	}
	return traps
}
