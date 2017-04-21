package geom

import "github.com/200sc/go-compgeo/search"

// Dot2D returns the dot product
// of two 2d elements.
func Dot2D(p1, p2 D2) float64 {
	return p1.X()*p2.X() + p1.Y()*p2.Y()
}

// Cross2D preforms the cross product on three points
// in two dimensions.
func Cross2D(a, b, c D2) float64 {
	return (b.X()-a.X())*(c.Y()-a.Y()) -
		(b.Y()-a.Y())*(c.X()-a.X())
}

// VertCross2D returns the cross product
// corrected for verticality-- b and c
// are organized such that left/rightness
// checks can be made.
func VertCross2D(a, b, c D2) float64 {
	cp := Cross2D(a, b, c)
	// If the first point of the line is above the second,
	// the cross product will return a negative value for left.
	if b.Y() > c.Y() {
		cp *= -1
	}
	return cp
}

// HzCross2D is equivalent to VertCross2D
// for horizontal queries.
func HzCross2D(a, b, c D2) float64 {
	cp := Cross2D(a, b, c)
	if b.X() > c.X() {
		cp *= -1
	}
	return cp
}

// IsColinear returns whether the cross product reports 0.
func IsColinear(a, b, c D2) bool {
	return Cross2D(a, b, c) == 0
}

// IsAbove returns whether a is above the line segment
// formed by (b->c)
func IsAbove(a, b, c D2) bool {
	return HzCross2D(a, b, c) > 0
}

// IsBelow is equivalent to !IsAbove && !IsColinear
func IsBelow(a, b, c D2) bool {
	return HzCross2D(a, b, c) < 0
}

// IsColinearOrAbove is equivalent to calling IsColinear || IsAbove
// without redoing the cross product calculation.
func IsColinearOrAbove(a, b, c D2) bool {
	return HzCross2D(a, b, c) >= 0
}

// IsColinearOrBelow is equivalent to calling IsColinear || IsBelow
// without redoing the cross product calculation
func IsColinearOrBelow(a, b, c D2) bool {
	return HzCross2D(a, b, c) <= 0
}

// IsLeftOf returns whether a is to the left of the line segment
// formed by (b->c)
func IsLeftOf(a, b, c D2) bool {
	return VertCross2D(a, b, c) > 0
}

// IsRightOf is equivalent to !IsLeftOf, except that both
// IsRightOf and IsLeftOf return false for Cross2D() == 0
func IsRightOf(a, b, c D2) bool {
	return VertCross2D(a, b, c) < 0
}

// IsColinearOrLeft is equivalent to calling IsColinear || IsLeft
// without redoing the cross product calculation.
func IsColinearOrLeft(a, b, c D2) bool {
	return VertCross2D(a, b, c) >= 0
}

// IsColinearOrRight is equivalent to calling IsColinear || IsRight
// without redoing the cross product calculation
func IsColinearOrRight(a, b, c D2) bool {
	return VertCross2D(a, b, c) <= 0
}

// VerticalCompare returns a search result
// representing whether this point is above
// equal or below the query edge.
func VerticalCompare(dp D2, e Spanning) search.CompareResult {
	p1 := e.At(0).(D2)
	p2 := e.At(1).(D2)
	if p1.X() < p2.X() {
		p1, p2 = p2, p1
	}
	s := Cross2D(p1, p2, dp)
	if s == 0 {
		return search.Equal
	} else if s < 0 {
		return search.Less
	}
	return search.Greater
}
