package kirkpatrick

import (
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/monotone"
	"github.com/200sc/go-compgeo/dcel/trapezoid"
)

//Triangulation method constant
type Method int

const (
	MONOTONE Method = iota
	TRAPEZOID
)

func TriangleTree(dc *dcel.DCEL, m Method) (dcel.LocatesPoints, error) {
	var tri *dcel.DCEL
	var mp map[*dcel.Face]*dcel.Face
	var err error
	switch m {
	case MONOTONE:
		// We need to wrap this dcel in some bounding polygon.
		// for our purposes we treat
		tri, mp, err = monotone.Triangulate(dc)
	case TRAPEZOID:
		// The trapezoidal map method requires that we add to our
		// dcel a wrapping square, so we satisfy kirkpatrick's
		// wrapping method
		tri, mp, _, err = trapezoid.TrapezoidalMap(dc)
		if err != nil {
			return nil, err
		}
		tri, mp, err = monotone.TriangulateSplit(tri, mp)
	}
	if err != nil {
		return nil, err
	}
	// ...
	return nil, nil
}
