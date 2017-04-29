package rtree

import (
	"fmt"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
	"github.com/sythe2o0/rtreego"
)

func DCELtoRtree(dc *dcel.DCEL) *rtreego.Rtree {
	tree := rtreego.NewTree(20, 40)

	for _, f := range dc.Faces {
		tree.Insert(&SpatialFace{f})
	}

	return tree
}

type SpatialFace struct {
	*dcel.Face
}

func (sf *SpatialFace) Bounds() *rtreego.Rect {
	span := sf.Face.Bounds()
	min := span.Left()
	max := span.Right()
	p := rtreego.Point{min.X(), min.Y(), min.Z()}
	rect, err := rtreego.NewRect(p, [3]float64{max.X(), max.Y(), max.Z()})
	if err != nil {
		fmt.Println("bounds error:", err)
	}
	return &rect
}

// SearchIntersect filters the output of rtree.SearchIntersect
// on plumb line contains for the tree's faces
func SearchIntersect(tree *rtreego.Rtree, p geom.D3) []*dcel.Face {
	pt := rtreego.Point{p.X(), p.Y(), p.Z()}
	rect, err := rtreego.NewRect(pt, [3]float64{0.1, 0.1, 0.1})
	if err != nil {
		fmt.Println("bounds error:", err)
	}
	spts := tree.SearchIntersect(&rect)
	out := make([]*dcel.Face, 0)
	for _, s := range spts {
		sf := s.(*SpatialFace)
		if sf.Contains(p) {
			out = append(out, sf.Face)
		}
	}
	return out
}
