package slab

import (
	"fmt"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search"
)

type face struct {
	*dcel.Face
}

func (f face) Equals(e search.Equalable) bool {
	switch f2 := e.(type) {
	case face:
		return f.Face == f2.Face
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
		compX, err := ce.FindSharedPoint(c.Edge, 0)
		if err != nil {
			fmt.Println("Edges share no point on x axis", ce, c)
		}
		p1, _ := ce.PointAt(0, compX)
		p2, _ := c.PointAt(0, compX)
		if p1[1] < p2[1] {
			return search.Less
		}
		return search.Greater
	}
	return ce.Edge.Compare(i)
}
