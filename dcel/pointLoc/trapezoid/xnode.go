package trapezoid

import (
	"image/color"

	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
	"github.com/200sc/go-compgeo/geom"
)

// NewX returns an X-Node at point P
func NewX(p geom.D3) *Node {
	return &Node{
		query:   xQuery,
		payload: p,
	}
}

func xQuery(fe geom.FullEdge, n *Node) []*Trapezoid {
	p := n.payload.(geom.Point)
	if visualize.VisualCh != nil {
		visualize.HighlightColor = color.RGBA{128, 128, 128, 128}
		visualize.DrawVerticalLine(p)
	}
	if geom.F64eq(fe.Left().X(), p.X()) {
		// If equal, go right.
		return n.right.Query(fe)
	} else if fe.Left().X() < p.X() {
		return n.left.Query(fe)
	}
	return n.right.Query(fe)
}
