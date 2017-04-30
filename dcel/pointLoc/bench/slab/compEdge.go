package slab

import (
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search"
)

type faces struct {
	f1, f2 *dcel.Face
}

func (fs faces) Equals(e search.Equalable) bool {
	switch fs2 := e.(type) {
	case faces:
		return fs2.f1 == fs.f1 && fs2.f2 == fs.f2
	}
	return false
}

type shellNode struct {
	k compEdge
	v search.Equalable
}

func (sn shellNode) Key() search.Comparable {
	return sn.k
}

func (sn shellNode) Val() search.Equalable {
	return sn.v
}

type compEdge struct {
	*dcel.Edge
}

func (ce compEdge) Compare(i interface{}) search.CompareResult {
	switch c := i.(type) {
	case compEdge:
		if ce.Edge == c.Edge {
			return search.Equal
		}

		if geom.F64eq(ce.X(), c.X()) && geom.F64eq(ce.Y(), c.Y()) &&
			geom.F64eq(ce.Twin.X(), c.Twin.X()) && geom.F64eq(ce.Twin.Y(), c.Twin.Y()) {
			return search.Equal
		}
		compX, _ := ce.FindSharedPoint(c.Edge, 0)
		p1, _ := ce.PointAt(0, compX)
		p2, _ := c.PointAt(0, compX)
		if p1[1] < p2[1] {
			return search.Less
		}
		return search.Greater
	}
	return ce.Edge.Compare(i)
}
