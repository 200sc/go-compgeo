package dcel

import (
	"math"

	"github.com/200sc/go-compgeo/printutil"
	"github.com/200sc/go-compgeo/search"
)

const (
	POINT_DIM = 3
)

// A Point is just a N-dimensional point.
type Point [POINT_DIM]float64

// NewPoint returns a Point initialized at the given position
func NewPoint(x, y, z float64) *Point {
	return &Point{x, y, z}
}

// String converts dp into a string.
func (dp Point) String() string {
	s := ""
	s += "("
	s += printutil.Stringf64(dp[0], dp[1], dp[2])
	s += ")"
	return s
}

// D returns the number of dimensions supported by
// a point.
func (dp Point) D() int {
	return POINT_DIM
}

// Val is equivalent to array access on a point.
func (dp Point) Val(d int) float64 {
	return dp[d]
}

// X :
// Get the value of this point on the x axis
func (dp Point) X() float64 {
	return dp[0]
}

// Y :
// Get the value of this point on the y axis
func (dp Point) Y() float64 {
	return dp[1]
}

// Z :
// Get the value of this point on the z axis
func (dp Point) Z() float64 {
	return dp[2]
}

// Eq returns whether two points are equivalent.
func (dp Point) Eq(p2 Dimensional) bool {
	if dp.D() != p2.D() {
		return false
	}
	for i := range dp {
		if dp[i] != p2.Val(i) {
			return false
		}
	}
	return true
}

// Mid2D returns the point in the middle of
// this point and p2.
func (dp Point) Mid2D(p2 D2) *Point {
	p3 := new(Point)
	for i := range dp {
		p3[i] = (dp[i] + p2.Val(i)) / 2
	}
	return p3
}

// VerticalCompare returns a search result
// representing whether this point is above
// equal or below the query edge.
func (dp Point) VerticalCompare(e *Edge) search.CompareResult {
	p1 := e.Origin.Point
	p2 := e.Twin.Origin.Point
	if p1[0] < p2[0] {
		p1, p2 = p2, p1
	}
	s := p1.Cross2D(p2, dp)
	if s == 0 {
		return search.Equal
	} else if s < 0 {
		return search.Less
	}
	return search.Greater
}

// Dot2D performs Dot multiplication on the two
// points, in a two-dimensional context.
func (dp Point) Dot2D(p2 *Point) float64 {
	return dp[0]*p2[0] + dp[1]*p2[1]
}

// Cross2D performs the Cross Product on the three
// points, in a two-dimensional context.
func (dp Point) Cross2D(p2, p3 D2) float64 {
	return Cross2D(dp, p2, p3)
}

// Lesser2D reports the lower point by y value,
// or by x value given equal y values. If the
// two points are equal the latter point is returned.
func (dp Point) Lesser2D(p2 D2) D2 {
	if dp[1] < p2.Val(1) {
		return dp
	} else if dp[1] > p2.Val(1) {
		return p2
	}
	if dp[0] < p2.Val(0) {
		return dp
	}
	return p2
}

// Greater2D reports the higher point by y value,
// or by x value given equal y values. If the
// two points are equal the latter point is returned.
func (dp Point) Greater2D(p2 D2) D2 {
	p3 := dp.Lesser2D(p2)
	if p3.Eq(dp) {
		return p2
	}
	return dp
}

// Magnitude2D reports the magnitude of the point
// interpreted as a vector in a two-dimensional context.
func (dp Point) Magnitude2D() float64 {
	return math.Sqrt((dp[0] * dp[0]) + (dp[1] * dp[1]))
}
