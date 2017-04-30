package rtree

import (
	"fmt"

	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
	"github.com/sythe2o0/rtreego"
)

func DCELtoRtree(dc *dcel.DCEL) *Rtree {
	tree := rtreego.NewTree(20, 40)

	for i := 1; i < len(dc.Faces); i++ {
		tree.Insert(&SpatialFace{dc.Faces[i]})
	}

	return &Rtree{tree}
}

type SpatialFace struct {
	*dcel.Face
}

func (sf *SpatialFace) Bounds() *rtreego.Rect {
	span := sf.Face.Bounds()
	min := span.Left()
	diff := span.Diff()
	p := rtreego.Point{min.X(), min.Y(), min.Z()}
	dist := [3]float64{diff.X(), diff.Y(), diff.Z()}
	for i, v := range dist {
		if v <= 0 {
			dist[i] = 0.1
		}
	}
	rect, err := rtreego.NewRect(p, dist)
	if err != nil {
		fmt.Println("bounds error:", err)
	}
	return &rect
}

type Rtree struct {
	*rtreego.Rtree
}

func (rt *Rtree) PointLocate(vs ...float64) (*dcel.Face, error) {
	if len(vs) < 2 {
		return nil, compgeo.InsufficientDimensionsError{}
	}
	pt := geom.NewPoint(vs[0], vs[1], 0)
	if len(vs) > 2 {
		pt = pt.Set(2, vs[2]).(geom.Point)
	}
	fs := SearchIntersect(rt.Rtree, pt)
	if len(fs) > 0 {
		return fs[0], nil
	}
	return nil, nil
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
