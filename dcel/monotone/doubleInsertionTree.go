package monotone

import (
	compgeo "github.com/200sc/go-compgeo"
	"github.com/200sc/go-compgeo/dcel"
	"github.com/200sc/go-compgeo/geom"
	"github.com/200sc/go-compgeo/search"
	"github.com/200sc/go-compgeo/search/tree"
)

// DoubleIntervalTree converts a monotonized f into a
// structure that can be pointlocated on to determine if
// a given point exists inside or outside the face.
// Because this is a monotone polygon we don't actually
// make an interval tree, we just make bsts. The intervals
// we use will be non-overlapping except at vertices
func NewDoubleIntervalTree(f *dcel.Face, dc *dcel.DCEL) (dcel.LocatesPoints, error) {
	var e *dcel.Edge
	for e = f.Outer; VertexType(e.Origin, dc) != START; e = e.Next {
	}
	// e is now the edge whose origin is the start vertex.
	leftTree := tree.New(tree.RedBlack)
	rightTree := tree.New(tree.RedBlack)
	st := e
	leftTree.Insert(interval{st})
	tree := leftTree
	for e := st.Next; e != st; e = e.Next {
		if VertexType(e.Origin, dc) == END {
			tree = rightTree
		}
		tree.Insert(interval{e})
	}
	return DblIntervalTree{leftTree, rightTree, f}, nil
}

type DblIntervalTree struct {
	leftTree, rightTree search.Dynamic
	f                   *dcel.Face
}

func (dit DblIntervalTree) PointLocate(vs ...float64) (*dcel.Face, error) {
	if len(vs) < 2 {
		return nil, compgeo.InsufficientDimensionsError{}
	}
	found, i1 := dit.leftTree.Search(yVal(vs[1]))
	if !found {
		return nil, nil
	}
	found2, i2 := dit.rightTree.Search(yVal(vs[1]))
	if !found2 {
		return nil, nil
	}
	e1 := i1.(interval).Edge
	e2 := i2.(interval).Edge
	pt := geom.NewPoint(vs[0], vs[1], 0)
	// Consider-- this could probably be a direct comparison
	// to pointAt instead of a cross product
	c1 := geom.VertCross2D(pt, e1.Origin, e1.Twin.Origin)
	c2 := geom.VertCross2D(pt, e2.Origin, e2.Twin.Origin)
	// Check that either c1 or c2 is 0, or otherwise
	// that the signs of c1 and c2 are different.
	// We don't actually care if the left and right trees are on the left
	// and right respectively. If we are on the same side of both trees,
	// we aren't in the polygon. If we are on different sides of both trees,
	// as we can't both be to the left of the (actual) left tree and right
	// of the (actual) right tree, we'll be in the polygon.
	if c1 == 0 || c2 == 0 {
		return dit.f, nil
	} else if c1*c2 < 0 {
		return dit.f, nil
	}
	return nil, nil
}

type interval struct {
	*dcel.Edge
}

func (i interval) Key() search.Comparable {
	return i
}

func (i interval) Val() search.Equalable {
	return i
}

func (i interval) Equals(e search.Equalable) bool {
	switch i2 := e.(type) {
	case interval:
		return i == i2
	}
	return false
}

func (i interval) Compare(s interface{}) search.CompareResult {
	switch i2 := s.(type) {
	case interval:
		if geom.F64eq(i.Origin.Y(), i2.Origin.Y()) {
			if geom.F64eq(i.Twin.Origin.Y(), i2.Twin.Origin.Y()) {
				return search.Equal
			} else if i.Twin.Origin.Y() < i2.Twin.Origin.Y() {
				return search.Less
			}
			return search.Greater
		} else if i.Origin.Y() < i2.Origin.Y() {
			return search.Less
		}
		return search.Greater
	case yVal:
		return i2.Compare(i)
	}
	return search.Invalid
}

type yVal float64

func (yv yVal) Compare(s interface{}) search.CompareResult {
	switch i := s.(type) {
	case interval:
		y := float64(yv)
		if geom.F64eq(y, i.Origin.Y()) ||
			geom.F64eq(y, i.Twin.Origin.Y()) {
			return search.Equal
		}
		c1 := y < i.Origin.Y()
		c2 := y > i.Twin.Origin.Y()
		c3 := y > i.Origin.Y()
		c4 := y < i.Twin.Origin.Y()
		if (c1 && c2) || (c3 && c4) {
			return search.Equal
		}
		if c1 && c4 {
			return search.Less
		}
		return search.Greater
	}
	return search.Invalid
}
