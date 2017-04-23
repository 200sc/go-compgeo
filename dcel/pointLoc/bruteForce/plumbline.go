package bruteForce

import (
	"fmt"
	"image/color"

	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc"
	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
	"github.com/200sc/go-compgeo/geom"
)

// The Plumb Line method is a name for a linear PIP check that
// shoots a ray out and checks how many times that ray intersects
// a polygon. The variation on a DCEL will iteratively perform
// plumb line on each face of the DCEL.
func PlumbLine(dc *dcel.DCEL) pointLoc.LocatesPoints {
	return &Iterator{dc}
}

type Iterator struct {
	*dcel.DCEL
}

func (i *Iterator) PointLocate(vs ...float64) (*dcel.Face, error) {
	if len(vs) < 2 {
		return nil, compgeo.InsufficientDimensionsError{}
	}
	p := geom.NewPoint(vs[0], vs[1], 0)
	containFn := Contains
	if visualize.VisualCh != nil {
		containFn = VisualizeContains
	}
	for j := 1; j < len(i.Faces); j++ {
		f := i.Faces[j]
		if containFn(f, p) {
			return f, nil
		}
	}
	return nil, nil
}

func Contains(f *dcel.Face, p geom.D2) bool {
	return f.Contains(p)
}

// Contains returns whether a point lies inside f.
// We cannot assume that f is convex, or anything
// besides some polygon. That leaves us with a rather
// complex form of PIP--
func VisualizeContains(f *dcel.Face, p geom.D2) bool {
	x := p.X()
	y := p.Y()
	contains := false
	bounds := f.Bounds()
	min := bounds.At(0).(geom.D2)
	max := bounds.At(1).(geom.D2)
	fmt.Println("Face bounds", bounds)
	visualize.HighlightColor = color.RGBA{0, 0, 255, 255}
	visualize.DrawFace(f)
	if x < min.Val(0) || x > max.Val(0) ||
		y < min.Val(1) || y > max.Val(1) {
		return contains
	}

	e1 := f.Outer.Prev
	e2 := f.Outer
	for {
		visualize.HighlightColor = color.RGBA{0, 0, 255, 255}
		visualize.DrawLine(e2.Origin, e1.Origin)
		if (e2.Y() > y) != (e1.Y() > y) {
			if x < (e1.X()-e2.X())*(y-e2.Y())/(e1.Y()-e2.Y())+e2.X() {
				visualize.HighlightColor = color.RGBA{0, 255, 0, 255}
				visualize.DrawLine(e2.Origin, e1.Origin)
				contains = !contains
			}
		}
		e1 = e1.Next
		e2 = e2.Next
		if e1 == f.Outer.Prev {
			break
		}
	}
	return contains
}
