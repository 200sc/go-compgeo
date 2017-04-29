package trapezoid

import (
	"fmt"
	"image/color"

	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
	"github.com/200sc/go-compgeo/geom"
)

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
	if visualize.VisualCh != nil {
		visualize.HighlightColor = color.RGBA{0, 0, 128, 128}
		visualize.DrawPoly(tr.toPhysics())
	}
	r := fe.Right()
	fmt.Println("rights:", tr.right, r.X())
	for tr != nil && r.X() > tr.right {
		// We perform this check here is it is less expensive
		// than the cross product in the latter case, even
		// though the latter case would suffice to do this.
		if tr.Neighbors[upright] == tr.Neighbors[botright] {
			if tr.Neighbors[upright] == nil {
				fmt.Println("Neighbors", tr.Neighbors)
			}
			tr = tr.Neighbors[botright]
			fmt.Println("New tr, botright and upright equal")
			fmt.Println(tr)
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
			visualize.HighlightColor = visualize.CheckFaceColor
			visualize.DrawPoly(tr.toPhysics())
			traps = append(traps, tr)
		}
	}
	return traps
}
