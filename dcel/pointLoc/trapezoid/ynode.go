package trapezoid

import (
	"fmt"
	"image/color"

	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
	"github.com/200sc/go-compgeo/geom"
)

// NewY returns a Y-Node at edge e
func NewY(e geom.FullEdge) *Node {
	return &Node{
		query:   yQuery,
		payload: e,
	}
}

func yQuery(fe geom.FullEdge, n *Node) []*Trapezoid {
	// This query asks if fe.Left() is above or below
	// yn.FullEdge.
	// If they are colinear, however, we need to check
	// which slope is larger. If fe is larger, we go above,
	// else we go below.
	yn := n.payload.(geom.FullEdge)
	if visualize.VisualCh != nil {
		visualize.HighlightColor = color.RGBA{128, 128, 128, 128}
		visualize.DrawLine(yn.Left(), yn.Right())
	}
	cp := geom.HzCross2D(fe.Left(), yn.Left(), yn.Right())
	fmt.Println("Y compare:", cp, fe, yn)
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
