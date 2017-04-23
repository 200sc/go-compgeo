package monotone

import (
	"fmt"

	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/dcel/visualize"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search"
)

// CompEdge describes structs to satisfy search interfaces
// so edges can be put as values and keys into a binary search
// tree and compared horizontally
//
// See also: slab/compEdge.go
// these structures, or at least the Compare function portion
// should probably be in dcel

type edgeNode struct {
	v *dcel.Edge
}

func (en edgeNode) Key() search.Comparable {
	return compEdge{en.v}
}

func (en edgeNode) Val() search.Equalable {
	return valEdge{en.v}
}

type valEdge struct {
	*dcel.Edge
}

func (ve valEdge) Equals(e search.Equalable) bool {
	switch ve2 := e.(type) {
	case valEdge:
		return ve.Edge == ve2.Edge
	}
	return false
}

// We need to have our keys be CompEdges so
// they are comparable within a certain y range.
type compEdge struct {
	*dcel.Edge
}

func (ce compEdge) Compare(i interface{}) search.CompareResult {
	switch c := i.(type) {
	case compEdge:
		if visualize.VisualCh != nil {
			visualize.DrawLine(ce.Edge.Origin, ce.Edge.Twin.Origin)
			visualize.DrawLine(c.Edge.Origin, c.Edge.Twin.Origin)
		}
		fmt.Println("Comparing", ce, c)
		if ce.Edge == c.Edge {
			fmt.Println("Equal1!")
			return search.Equal
		}

		if geom.F64eq(ce.X(), c.X()) && geom.F64eq(ce.Y(), c.Y()) &&
			geom.F64eq(ce.Twin.X(), c.Twin.X()) && geom.F64eq(ce.Twin.Y(), c.Twin.Y()) {
			fmt.Println("Equal2!")
			return search.Equal
		}
		y, err := ce.FindSharedPoint(c.Edge, 1)
		if err != nil {
			fmt.Println("Edges share no y point")
		}
		p1, _ := ce.PointAt(1, y)
		p2, _ := c.PointAt(1, y)
		if p1[0] < p2[0] {
			fmt.Println("Less!")
			return search.Less
		}
		fmt.Println("Greater!")
		return search.Greater
	}
	return ce.Edge.Compare(i)
}
