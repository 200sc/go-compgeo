package slab

import (
	"fmt"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/pointLoc/visualize"
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
		visualize.DrawLine(ce.Edge.Origin, ce.Edge.Twin.Origin)
		visualize.DrawLine(c.Edge.Origin, c.Edge.Twin.Origin)
		if ce.Edge == c.Edge {
			return search.Equal
		}

		if geom.F64eq(ce.X(), c.X()) && geom.F64eq(ce.Y(), c.Y()) &&
			geom.F64eq(ce.Twin.X(), c.Twin.X()) && geom.F64eq(ce.Twin.Y(), c.Twin.Y()) {
			return search.Equal
		}
		compX, err := ce.FindSharedPoint(c.Edge, 0)
		if err != nil {
			fmt.Println("Edges share no point on x axis")
		}
		p1, err := ce.PointAt(0, compX)
		if err != nil {
			fmt.Println(ce, " Did not have a value at ", compX)
		}
		p2, err := c.PointAt(0, compX)
		if err != nil {
			fmt.Println(c, " Did not have a value at ", compX)
		}
		if p1[1] < p2[1] {
			return search.Less
		}
		if p1[1] > p2[1] {
			return search.Greater
		}
		fmt.Println("Something went wrong")
		return search.Greater
	}
	return ce.Edge.Compare(i)
}
