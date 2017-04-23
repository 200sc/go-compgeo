package trapezoid

import (
	"fmt"
	"image/color"

	"github.com/200sc/go-compgeo/dcel/visualize"
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
	if fe.Left().X() < p.X() {
		fmt.Println("X compare:", fe.Left().X(), p.X(), true)
		return n.left.Query(fe)
	}
	fmt.Println("X compare:", fe.Left().X(), p.X(), false)
	return n.right.Query(fe)
}
