// bruteForce implements some simple brute force point location methods

package bruteForce

import (
	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc"
	"github.com/200sc/go-compgeo/geom"
)

// PlumbLine method is a name for a linear PIP check that
// shoots a ray out and checks how many times that ray intersects
// a polygon. The variation on a DCEL will iteratively perform
// plumb line on each face of the DCEL.
func PlumbLine(dc *dcel.DCEL) pointLoc.LocatesPoints {
	return &Iterator{dc}
}

// Iterator is a simple dcel wrapper for the following pointLocate method
type Iterator struct {
	*dcel.DCEL
}

// PointLocate on an iterator performs plumb line on each
// of a DCEL's faces in order.
func (i *Iterator) PointLocate(vs ...float64) (*dcel.Face, error) {
	if len(vs) < 2 {
		return nil, compgeo.InsufficientDimensionsError{}
	}
	p := geom.NewPoint(vs[0], vs[1], 0)
	for j := 1; j < len(i.Faces); j++ {
		f := i.Faces[j]
		if f.Contains(p) {
			return f, nil
		}
	}
	return nil, nil
}
